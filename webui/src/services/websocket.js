/**
 * WebSocket Service for Real-time Communication
 * Handles real-time messaging, presence updates, and live notifications
 */

class WebSocketService {
	constructor() {
		this.ws = null;
		this.isConnected = false;
		this.reconnectAttempts = 0;
		this.maxReconnectAttempts = 5;
		this.reconnectInterval = 1000; // Start with 1 second
		this.listeners = new Map();
		this.heartbeatInterval = null;
		this.userId = null;
	}

	/**
	 * Connect to WebSocket server
	 * @param {string} userId - Current user ID
	 */
	connect(userId) {
		if (this.isConnected) {
			console.log('WebSocket already connected');
			return;
		}

		this.userId = userId;
		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		const wsUrl = `${protocol}//${window.location.host}/ws?user_id=${userId}`;

		try {
			this.ws = new WebSocket(wsUrl);
			this.setupEventHandlers();
		} catch (error) {
			console.error('Failed to create WebSocket connection:', error);
			this.scheduleReconnect();
		}
	}

	/**
	 * Setup WebSocket event handlers
	 */
	setupEventHandlers() {
		this.ws.onopen = () => {
			console.log('WebSocket connected');
			this.isConnected = true;
			this.reconnectAttempts = 0;
			this.reconnectInterval = 1000;
			this.startHeartbeat();
			this.emit('connected');
		};

		this.ws.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data);
				this.handleMessage(data);
			} catch (error) {
				console.error('Failed to parse WebSocket message:', error);
			}
		};

		this.ws.onclose = (event) => {
			console.log('WebSocket closed:', event.code, event.reason);
			this.isConnected = false;
			this.stopHeartbeat();
			this.emit('disconnected');

			// Attempt to reconnect unless it was a clean close
			if (event.code !== 1000) {
				this.scheduleReconnect();
			}
		};

		this.ws.onerror = (error) => {
			console.error('WebSocket error:', error);
			this.emit('error', error);
		};
	}

	/**
	 * Handle incoming WebSocket messages
	 * @param {Object} data - Message data
	 */
	handleMessage(data) {
		const { type, payload } = data;

		switch (type) {
			case 'message':
				this.emit('message', payload);
				break;
			case 'message_deleted':
				this.emit('messageDeleted', payload);
				break;
			case 'reaction_added':
			case 'reaction_removed':
				this.emit('reactionChanged', payload);
				break;
			// Comment events removed
			case 'user_online':
				this.emit('userOnline', payload);
				break;
			case 'user_offline':
				this.emit('userOffline', payload);
				break;
			case 'conversation_updated':
				this.emit('conversationUpdated', payload);
				break;
			case 'typing_start':
				this.emit('typingStart', payload);
				break;
			case 'typing_stop':
				this.emit('typingStop', payload);
				break;
			case 'pong':
				// Heartbeat response - connection is alive
				break;
			default:
				console.warn('Unknown WebSocket message type:', type);
		}
	}

	/**
	 * Send a message through WebSocket
	 * @param {string} type - Message type
	 * @param {Object} payload - Message payload
	 */
	send(type, payload = {}) {
		if (!this.isConnected || !this.ws) {
			console.warn('WebSocket not connected, cannot send message');
			return false;
		}

		try {
			const message = JSON.stringify({ type, payload });
			this.ws.send(message);
			return true;
		} catch (error) {
			console.error('Failed to send WebSocket message:', error);
			return false;
		}
	}

	/**
	 * Subscribe to WebSocket events
	 * @param {string} event - Event name
	 * @param {Function} callback - Event callback
	 */
	on(event, callback) {
		if (!this.listeners.has(event)) {
			this.listeners.set(event, []);
		}
		this.listeners.get(event).push(callback);
	}

	/**
	 * Unsubscribe from WebSocket events
	 * @param {string} event - Event name
	 * @param {Function} callback - Event callback to remove
	 */
	off(event, callback) {
		if (!this.listeners.has(event)) return;

		const callbacks = this.listeners.get(event);
		const index = callbacks.indexOf(callback);
		if (index > -1) {
			callbacks.splice(index, 1);
		}
	}

	/**
	 * Emit event to all listeners
	 * @param {string} event - Event name
	 * @param {*} data - Event data
	 */
	emit(event, data) {
		if (!this.listeners.has(event)) return;

		const callbacks = this.listeners.get(event);
		callbacks.forEach((callback) => {
			try {
				callback(data);
			} catch (error) {
				console.error(
					`Error in WebSocket event callback for ${event}:`,
					error
				);
			}
		});
	}

	/**
	 * Start heartbeat to keep connection alive
	 */
	startHeartbeat() {
		this.heartbeatInterval = setInterval(() => {
			if (this.isConnected) {
				this.send('ping');
			}
		}, 30000); // Send ping every 30 seconds
	}

	/**
	 * Stop heartbeat
	 */
	stopHeartbeat() {
		if (this.heartbeatInterval) {
			clearInterval(this.heartbeatInterval);
			this.heartbeatInterval = null;
		}
	}

	/**
	 * Schedule reconnection attempt
	 */
	scheduleReconnect() {
		if (this.reconnectAttempts >= this.maxReconnectAttempts) {
			console.error('Max reconnection attempts reached');
			this.emit('maxReconnectAttemptsReached');
			return;
		}

		this.reconnectAttempts++;
		console.log(
			`Scheduling reconnect attempt ${this.reconnectAttempts} in ${this.reconnectInterval}ms`
		);

		setTimeout(() => {
			if (!this.isConnected) {
				console.log(`Reconnection attempt ${this.reconnectAttempts}`);
				this.connect(this.userId);
			}
		}, this.reconnectInterval);

		// Exponential backoff: double the interval for next attempt, max 30 seconds
		this.reconnectInterval = Math.min(this.reconnectInterval * 2, 30000);
	}

	/**
	 * Manually disconnect WebSocket
	 */
	disconnect() {
		if (this.ws) {
			this.ws.close(1000, 'Manual disconnect');
			this.ws = null;
		}
		this.isConnected = false;
		this.stopHeartbeat();
	}

	/**
	 * Send typing indicator
	 * @param {string} conversationId - Conversation ID
	 * @param {boolean} isTyping - Whether user is typing
	 */
	sendTypingIndicator(conversationId, isTyping) {
		this.send(isTyping ? 'typing_start' : 'typing_stop', {
			conversation_id: conversationId,
			user_id: this.userId,
		});
	}

	/**
	 * Join a conversation for real-time updates
	 * @param {string} conversationId - Conversation ID to join
	 */
	joinConversation(conversationId) {
		this.send('join_conversation', {
			conversation_id: conversationId,
		});
	}

	/**
	 * Leave a conversation
	 * @param {string} conversationId - Conversation ID to leave
	 */
	leaveConversation(conversationId) {
		this.send('leave_conversation', {
			conversation_id: conversationId,
		});
	}

	/**
	 * Get connection status
	 * @returns {boolean} Whether WebSocket is connected
	 */
	getConnectionStatus() {
		return this.isConnected;
	}
}

// Create and export singleton instance
export const webSocketService = new WebSocketService();
export default webSocketService;
