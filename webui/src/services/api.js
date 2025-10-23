import axios from './axios.js';

/**
 * WASAText API Service
 * Comprehensive API service implementing all OpenAPI endpoints
 */

// ============ AUTHENTICATION ============
export const auth = {
	/**
	 * User login/registration
	 * @param {string} name - Username (3-16 characters)
	 * @returns {Promise<{identifier: string}>}
	 */
	async login(name) {
		const response = await axios.post('/session', { name });
		return response.data;
	},
};

// ============ USERS ============
export const users = {
	/**
	 * List all users in the system
	 * @returns {Promise<User[]>}
	 */
	async listAll() {
		const response = await axios.get('/users');
		return response.data;
	},

	/**
	 * Update user's username
	 * @param {string} userId - User UUID
	 * @param {string} username - New username (3-16 characters)
	 * @returns {Promise<{message: string}>}
	 */
	async updateUsername(userId, username) {
		// Backend expects a raw JSON string for username, not an object
		const response = await axios.put(
			`/users/${userId}`,
			JSON.stringify(username),
			{ headers: { 'Content-Type': 'application/json' } }
		);
		return response.data;
	},

	/**
	 * Update user's profile picture
	 * @param {string} userId - User UUID
	 * @param {string} photoUrl - URL to new profile picture
	 * @returns {Promise<{message: string}>}
	 */
	async updatePhoto(userId, photoUrl) {
		const response = await axios.put(`/users/${userId}/photo`, photoUrl, {
			headers: { 'Content-Type': 'application/json' },
		});
		return response.data;
	},

	/**
	 * Get user context information
	 * @param {string} userId - User UUID
	 * @returns {Promise<object>}
	 */
	async getContext(userId) {
		const response = await axios.get(`/users/${userId}/context`);
		return response.data;
	},
};

// ============ CONVERSATIONS ============
export const conversations = {
	/**
	 * Get all conversations for a user
	 * @param {string} userId - User UUID
	 * @returns {Promise<Conversation[]>}
	 */
	async getUserConversations(userId) {
		const response = await axios.get(`/users/${userId}/conversations`);
		return response.data;
	},

	/**
	 * Create a new conversation
	 * @param {string} userId - User UUID
	 * @param {string[]} participants - Array of participant UUIDs
	 * @param {string} [name] - Optional conversation name for groups
	 * @returns {Promise<Conversation>}
	 */
	async create(userId, participants, name) {
		const data = { participants };
		if (name) data.name = name;

		const response = await axios.post(
			`/users/${userId}/conversations`,
			data
		);
		return response.data;
	},

	/**
	 * Get specific conversation details
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Conversation UUID
	 * @returns {Promise<Conversation>}
	 */
	async getDetails(userId, conversationId) {
		const response = await axios.get(
			`/users/${userId}/conversations/${conversationId}`
		);
		return response.data;
	},

	/**
	 * Add member to group conversation
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Conversation UUID
	 * @param {string} newUserId - UUID of user to add
	 * @returns {Promise<{message: string}>}
	 */
	async addMember(userId, conversationId, newUserName) {
		const response = await axios.post(
			`/users/${userId}/conversations/${conversationId}/members`,
			{ name: newUserName }
		);
		return response.data;
	},

	/**
	 * Leave a group conversation
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Conversation UUID
	 * @returns {Promise<{message: string}>}
	 */
	async leaveGroup(userId, conversationId) {
		const response = await axios.delete(
			`/users/${userId}/conversations/${conversationId}/members`
		);
		return response.data;
	},

	/**
	 * Set group conversation name
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Conversation UUID
	 * @param {string} name - New group name
	 * @returns {Promise<{message: string}>}
	 */
	async setGroupName(userId, conversationId, name) {
		const response = await axios.put(
			`/users/${userId}/conversations/${conversationId}/name`,
			name,
			{ headers: { 'Content-Type': 'application/json' } }
		);
		return response.data;
	},

	/**
	 * Set group conversation photo
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Conversation UUID
	 * @param {string} photoUrl - URL to new group photo
	 * @returns {Promise<{message: string}>}
	 */
	async setGroupPhoto(userId, conversationId, photoUrl) {
		const response = await axios.put(
			`/users/${userId}/conversations/${conversationId}/photo`,
			photoUrl,
			{ headers: { 'Content-Type': 'application/json' } }
		);
		return response.data;
	},

	/**
	 * Get all conversations in the system (admin function)
	 * @returns {Promise<Conversation[]>}
	 */
	// Removed admin-only function getAll()
};

// ============ MESSAGES ============
export const messages = {
	/**
	 * Get all messages from a conversation
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Conversation UUID
	 * @returns {Promise<Message[]>}
	 */
	async getConversationMessages(userId, conversationId) {
		const response = await axios.get(
			`/users/${userId}/conversations/${conversationId}/messages`
		);
		return response.data;
	},

	/**
	 * Send a message to a conversation
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Conversation UUID
	 * @param {string} [content] - Text content
	 * @param {string} [imageUrl] - Base64 image data URL
	 * @returns {Promise<string>} Message UUID
	 */
	async send(userId, conversationId, content, imageUrl) {
		const data = {};
		if (content) data.content = content;
		if (imageUrl) data.imageUrl = imageUrl;

		const response = await axios.post(
			`/users/${userId}/conversations/${conversationId}/messages`,
			data
		);
		return response.data;
	},

	/**
	 * Delete a message (only by sender)
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Conversation UUID
	 * @param {string} messageId - Message UUID
	 * @returns {Promise<{message: string}>}
	 */
	async delete(userId, conversationId, messageId) {
		const response = await axios.delete(
			`/users/${userId}/conversations/${conversationId}/messages/${messageId}`
		);
		return response.data;
	},

	/**
	 * Forward a message to another conversation
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Source conversation UUID
	 * @param {string} messageId - Message UUID to forward
	 * @param {string} targetConversationId - Target conversation UUID
	 * @returns {Promise<Conversation>}
	 */
	async forward(userId, conversationId, messageId, targetConversationId) {
		const response = await axios.post(
			`/users/${userId}/conversations/${conversationId}/messages/${messageId}/forward`,
			{ content: targetConversationId }
		);
		return response.data;
	},
};

// ============ REACTIONS ============
// ============ COMMENTS (emoji) ============
export const comments = {
	/**
	 * Toggle a reaction on a message
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Conversation UUID
	 * @param {string} messageId - Message UUID
	 * @param {string} emoji - Emoji to toggle
	 * @returns {Promise<{message: string}>}
	 */
	async toggle(userId, conversationId, messageId, emoji) {
		const response = await axios.post(
			`/users/${userId}/conversations/${conversationId}/messages/${messageId}/comments`,
			{ emoji }
		);
		return response.data;
	},

	/**
	 * Remove a specific reaction from a message
	 * @param {string} userId - User UUID
	 * @param {string} conversationId - Conversation UUID
	 * @param {string} messageId - Message UUID
	 * @param {string} emoji - Emoji to remove
	 * @returns {Promise<{message: string}>}
	 */
	async remove(userId, conversationId, messageId, emoji) {
		const response = await axios.delete(
			`/users/${userId}/conversations/${conversationId}/messages/${messageId}/comments/${emoji}`
		);
		return response.data;
	},
};

// ============ CONTACTS ============
export const contacts = {
	/**
	 * List all contacts for a user
	 * @param {string} userId - User UUID
	 * @returns {Promise<Contact[]>}
	 */
	async list(userId) {
		const response = await axios.get(`/users/${userId}/contacts`);
		return response.data;
	},

	/**
	 * Add a new contact
	 * @param {string} userId - User UUID
	 * @param {string} contactUserId - UUID of user to add as contact
	 * @returns {Promise<Contact>}
	 */
	async add(userId, contactUserId) {
		const response = await axios.post(`/users/${userId}/contacts`, {
			contactUserId,
		});
		return response.data;
	},

	/**
	 * Remove a contact
	 * @param {string} userId - User UUID
	 * @param {string} contactId - Contact UUID
	 * @returns {Promise<{message: string}>}
	 */
	async remove(userId, contactId) {
		const response = await axios.delete(
			`/users/${userId}/contacts/${contactId}`
		);
		return response.data;
	},
};

// ============ WEBSOCKET ============
export const websocket = {
	/**
	 * Create WebSocket connection for real-time messaging
	 * @returns {WebSocket}
	 */
	connect() {
		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		const host = window.location.host;
		const wsUrl = `${protocol}//${host}/ws`;

		return new WebSocket(wsUrl);
	},
};

// ============ HEALTH ============
export const health = {
	/**
	 * Get API root information and status
	 * @returns {Promise<object>} API information
	 */
	async getInfo() {
		const response = await axios.get('/');
		return response.data;
	},

	/**
	 * Check server health status
	 * @returns {Promise<string>} Health status
	 */
	async check() {
		const response = await axios.get('/liveness');
		return response.data;
	},

	/**
	 * Get recent system logs
	 * @returns {Promise<object>} Logs data
	 */
	// Removed admin endpoints: getLogs, checkHealth, getStats, getOnlineUsers
};

// Default export with all services
// Create the API service object
const apiService = {
	auth,
	users,
	conversations,
	messages,
	comments,
	contacts,
	websocket,
	health,
};

// Export both named and default
export { apiService };
export default apiService;
