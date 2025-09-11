<template>
	<div class="sidebar-content">
		<div class="sidebar-user">
			<img
				v-if="userPicture"
				:src="userPicture"
				alt="Profile"
				class="avatar"
			/>
			<div v-else class="avatar bg-secondary"></div>
			<div class="user-info">
				<div class="user-name">
					{{ username }}
					<button
						class="edit-profile-btn"
						@click="showEditProfile = true"
						title="Edit Profile"
					>
						✏️
					</button>
				</div>
				<div class="user-status">Online</div>
			</div>
			<div class="user-actions">
				<button
					class="theme-toggle-btn"
					@click="toggleTheme"
					:title="
						isDarkTheme
							? 'Switch to light theme'
							: 'Switch to dark theme'
					"
				>
					{{ isDarkTheme ? '☀️' : '🌙' }}
				</button>
				<button class="logout-btn" @click="handleLogout" title="Logout">
					🚪
				</button>
			</div>
		</div>
		<div class="sidebar-search">
			<input
				type="text"
				placeholder="Search or start new chat"
				v-model="searchQuery"
				@input="filterChats"
			/>
		</div>
		<div class="sidebar-actions">
			<button
				class="btn btn-primary btn-sm w-100"
				data-bs-toggle="modal"
				data-bs-target="#createConversationModal"
			>
				+ New Chat or Group
			</button>
		</div>
		<div class="sidebar-chats">
			<div v-if="loading" class="sidebar-loading">Loading chats...</div>
			<div v-else-if="error" class="sidebar-error">{{ error }}</div>
			<div
				v-else-if="filteredChats.length === 0 && searchQuery.trim()"
				class="sidebar-empty"
			>
				No chats match your search.
			</div>
			<div v-else-if="filteredChats.length === 0" class="sidebar-empty">
				No conversations
			</div>
			<div
				v-else
				class="sidebar-chat"
				v-for="chat in filteredChats"
				:key="chat.id"
				:class="{ selected: chat.id === selectedChatId }"
				@click="$emit('select-chat', chat.id)"
			>
				<img class="chat-avatar" :src="getChatAvatar(chat)" alt="" />
				<div class="chat-info">
					<div class="chat-name-row">
						<span class="chat-name">
							{{ getChatDisplayName(chat) }}
						</span>
						<button
							class="delete-chat-btn"
							@click.stop="deleteChat(chat)"
							title="Delete conversation"
						>
							🗑️
						</button>
					</div>
					<div class="chat-last-row">
						<span class="last-message" v-if="chat.lastMessage">
							<span
								class="last-sender"
								v-if="chat.lastMessage.senderId !== userId"
							>
								{{ chat.lastMessage.senderUsername }}:
							</span>
							<span
								v-if="chat.lastMessage.imageUrl"
								class="image-indicator"
								>📷
							</span>
							{{
								getLastMessagePreview(
									chat.lastMessage.text,
									chat.lastMessage.imageUrl
								)
							}}
						</span>
						<span class="no-messages" v-else>No messages yet</span>
						<span class="unread-count" v-if="chat.unreadCount > 0">
							{{ chat.unreadCount }}
						</span>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Edit Profile Modal -->
	<div
		v-if="showEditProfile"
		class="modal-overlay"
		@click="showEditProfile = false"
	>
		<div class="edit-profile-modal" @click.stop>
			<div class="modal-header">
				<h3>Edit Profile</h3>
				<button class="close-btn" @click="showEditProfile = false">
					×
				</button>
			</div>
			<div class="modal-body">
				<div class="form-group">
					<label for="newUsername">Username</label>
					<input
						id="newUsername"
						type="text"
						v-model="newUsername"
						:placeholder="username"
						class="form-control"
						maxlength="16"
						@keyup.enter="updateUsername"
					/>
					<small class="form-text text-muted">
						3-16 characters allowed
					</small>
				</div>
				<div v-if="editError" class="alert alert-danger">
					{{ editError }}
				</div>
			</div>
			<div class="modal-footer">
				<button
					class="btn btn-secondary"
					@click="showEditProfile = false"
				>
					Cancel
				</button>
				<button
					class="btn btn-primary"
					@click="updateUsername"
					:disabled="isUpdating || newUsername.length < 3"
				>
					{{ isUpdating ? 'Updating...' : 'Save' }}
				</button>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import axios from '../services/axios.js';
// Returns the avatar for a chat: for 1:1 chats, show the other participant's picture; for groups, show the chat picture
function getChatAvatar(chat) {
	const defaultAvatar =
		'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPGNpcmNsZSBjeD0iMjAiIGN5PSIyMCIgcj0iMjAiIGZpbGw9IiNlNWU3ZWIiLz4KPHN2ZyB4PSI4IiB5PSI4IiB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTEyIDEyYzIuMjEgMCA0LTEuNzkgNC00cy0xLjc5LTQtNC00LTQgMS43OS00IDQgMS43OSA0IDQgNHptMCAyYy0yLjY3IDAtOCAxLjM0LTggNHYyaDE2di0yYzAtMi42Ni01LjMzLTQtOC00eiIgZmlsbD0iIzlCA0E4Ii8+Cjwvc3ZnPgo8L3N2Zz4K';

	// Check if chat and participants exist
	if (!chat || !chat.participants || !Array.isArray(chat.participants)) {
		return defaultAvatar;
	}

	// Group chat: use chat.picture
	if (chat.participants.length > 2) {
		return chat.picture || defaultAvatar;
	}
	// 1:1 chat: use the other participant's picture
	if (chat.participants.length === 2) {
		const other = chat.participants.find(
			(p) => p && p.id && p.id !== props.userId
		);
		return other && other.picture
			? other.picture
			: chat.picture || defaultAvatar;
	}
	// Fallback
	return chat.picture || defaultAvatar;
}
// Returns the display name for a chat: for 1:1 chats, show the other participant's username; for groups, show the chat name
function getChatDisplayName(chat) {
	// Check if chat exists
	if (!chat) {
		return 'Unknown Chat';
	}

	// Check if participants exist
	if (!chat.participants || !Array.isArray(chat.participants)) {
		return chat.name || chat.id || 'Unknown Chat';
	}

	// If group chat (more than 2 participants)
	if (chat.participants.length > 2) {
		if (chat.name && chat.name.trim() !== '') {
			return chat.name;
		}
		// Fallback: join all usernames except the current user
		const otherUsers = chat.participants
			.filter((p) => p && p.id && p.id !== props.userId)
			.map((p) => p.username || 'Unknown')
			.filter(Boolean);
		return otherUsers.length > 0
			? otherUsers.join(', ')
			: chat.name || chat.id || 'Group Chat';
	}
	// For 1:1 chat, show the other participant's username
	if (chat.participants.length === 2) {
		// Find the participant who is not the current user
		const other = chat.participants.find(
			(p) => p && p.id && p.id !== props.userId
		);
		return other && other.username
			? other.username
			: chat.name || chat.id || 'Unknown User';
	}
	// Fallback
	return chat.name || chat.id || 'Unknown Chat';
}

// Creates a preview of the last message (truncated)
function getLastMessagePreview(text, imageUrl) {
	if (imageUrl && !text) {
		return 'Image';
	}
	if (!text) return '';
	const maxLength = 50;
	return text.length > maxLength
		? text.substring(0, maxLength) + '...'
		: text;
}

// Sort chats by last message time (newest first)
function sortChatsByLastMessage(chats) {
	return chats.sort((a, b) => {
		const aTime = a.lastMessageTime ? parseInt(a.lastMessageTime) : 0;
		const bTime = b.lastMessageTime ? parseInt(b.lastMessageTime) : 0;
		return bTime - aTime; // Descending order (newest first)
	});
}

// Update a specific chat with new message and reorder
function updateChatWithNewMessage(conversationId, newMessage) {
	const chatIndex = chats.value.findIndex(
		(chat) => chat.id === conversationId
	);
	if (chatIndex !== -1) {
		// Update the chat's last message
		chats.value[chatIndex].lastMessage = {
			id: newMessage.id,
			senderId: newMessage.senderId,
			text: newMessage.text || newMessage.content,
			senderUsername: newMessage.senderUsername,
		};
		// Update timestamp (use current time as approximation)
		chats.value[chatIndex].lastMessageTime = Date.now().toString();

		// Resort all chats
		const sortedChats = sortChatsByLastMessage([...chats.value]);
		chats.value = sortedChats;

		// Update filtered chats if there's a search query
		if (searchQuery.value.trim()) {
			filterChats();
		} else {
			filteredChats.value = sortedChats;
		}

		// Emit updated chats to parent
		emit('chats-loaded', chats.value);
	}
}

// Move a chat to the top when a new message is sent/received
function moveChartToTop(conversationId) {
	const chatIndex = chats.value.findIndex(
		(chat) => chat.id === conversationId
	);
	if (chatIndex !== -1 && chatIndex !== 0) {
		// Remove chat from current position and add to top
		const chat = chats.value.splice(chatIndex, 1)[0];
		chat.lastMessageTime = Date.now().toString();
		chats.value.unshift(chat);

		// Update filtered chats
		if (searchQuery.value.trim()) {
			filterChats();
		} else {
			filteredChats.value = [...chats.value];
		}

		// Emit updated chats to parent
		emit('chats-loaded', chats.value);
	}
}

const emit = defineEmits([
	'select-chat',
	'chats-loaded',
	'chat-deleted',
	'logout',
	'username-updated',
]);
const props = defineProps({
	userId: {
		type: String,
		required: true,
	},
	username: {
		type: String,
		required: true,
	},
	userPicture: {
		type: String,
		required: false,
		default: '',
	},
	selectedChatId: {
		type: [String, null],
		required: false,
	},
});

const chats = ref([]);
const filteredChats = ref([]);
const searchQuery = ref('');
const loading = ref(false);
const error = ref(null);
const isDarkTheme = ref(localStorage.getItem('darkTheme') === 'true');

// Profile editing variables
const showEditProfile = ref(false);
const newUsername = ref('');
const isUpdating = ref(false);
const editError = ref('');

async function fetchChats() {
	loading.value = true;
	error.value = null;
	try {
		const res = await axios.get(`/users/${props.userId}/conversations`);
		// Sort chats by last message time
		const sortedChats = sortChatsByLastMessage(res.data);
		chats.value = sortedChats;
		filteredChats.value = sortedChats; // Initialize filtered chats
		// Emit chats to parent so App.vue can sync selectedChat
		emit('chats-loaded', chats.value);
	} catch (err) {
		error.value = err.message || 'Failed to load chats';
		chats.value = [];
		filteredChats.value = [];
		emit('chats-loaded', []);
	} finally {
		loading.value = false;
	}
}

// Filter chats based on search query
function filterChats() {
	if (!searchQuery.value.trim()) {
		filteredChats.value = chats.value;
		return;
	}

	const query = searchQuery.value.toLowerCase();
	filteredChats.value = chats.value.filter((chat) => {
		// Search in conversation name
		if (chat.name && chat.name.toLowerCase().includes(query)) {
			return true;
		}

		// Search in participant names
		if (chat.participants) {
			return chat.participants.some(
				(participant) =>
					participant.username &&
					participant.username.toLowerCase().includes(query)
			);
		}

		return false;
	});
}

// Delete a chat/conversation
async function deleteChat(chat) {
	if (
		!confirm(
			`Are you sure you want to delete the conversation "${getChatDisplayName(
				chat
			)}"? This action cannot be undone.`
		)
	) {
		return;
	}

	try {
		// Use the correct leave group endpoint
		const deleteUrl = `/users/${props.userId}/conversations/${chat.id}/members`;
		console.log('Attempting to delete conversation with URL:', deleteUrl);
		console.log('User ID:', props.userId);
		console.log('Chat ID:', chat.id);

		await axios.delete(deleteUrl);

		console.log('Successfully deleted conversation');

		// Refresh the chats list
		await fetchChats();

		// If this was the selected chat, clear the selection
		emit('chat-deleted', chat.id);
	} catch (error) {
		console.error('Failed to delete conversation:', error);
		console.error('Error details:', error.response);

		let errorMessage = 'Failed to delete conversation. ';
		if (error.response) {
			if (error.response.status === 404) {
				errorMessage +=
					'Conversation not found or you are not a member.';
			} else if (error.response.status === 403) {
				errorMessage +=
					'You do not have permission to leave this conversation.';
			} else {
				errorMessage += `Server error: ${error.response.status}`;
			}
		} else {
			errorMessage += 'Please check your connection and try again.';
		}

		alert(errorMessage);
	}
}

// Toggle dark theme
function toggleTheme() {
	isDarkTheme.value = !isDarkTheme.value;
	localStorage.setItem('darkTheme', isDarkTheme.value.toString());
	applyTheme();
}

// Handle logout
function handleLogout() {
	emit('logout');
}

// Apply theme to document
function applyTheme() {
	if (isDarkTheme.value) {
		document.documentElement.classList.add('dark-theme');
	} else {
		document.documentElement.classList.remove('dark-theme');
	}
}

// Update username function
async function updateUsername() {
	if (newUsername.value.length < 3 || newUsername.value.length > 16) {
		editError.value = 'Username must be between 3 and 16 characters';
		return;
	}

	if (newUsername.value === props.username) {
		showEditProfile.value = false;
		return;
	}

	isUpdating.value = true;
	editError.value = '';

	try {
		const response = await axios.put(
			`/users/${props.userId}`,
			newUsername.value,
			{
				headers: {
					'Content-Type': 'application/json',
				},
			}
		);

		// Emit event to update username in parent component
		emit('username-updated', response.data.username);

		// Update localStorage
		localStorage.setItem('currentUsername', response.data.username);

		showEditProfile.value = false;
		newUsername.value = '';
	} catch (error) {
		console.error('Failed to update username:', error);
		if (error.response && error.response.status === 400) {
			editError.value = error.response.data || 'Username already in use';
		} else {
			editError.value = 'Failed to update username. Please try again.';
		}
	} finally {
		isUpdating.value = false;
	}
}

onMounted(() => {
	fetchChats();
	applyTheme(); // Apply theme on component mount
});
watch(() => props.userId, fetchChats);

// Expose methods for parent component
defineExpose({
	refreshChats: fetchChats,
	updateChatWithNewMessage,
	moveChartToTop,
});
</script>

<style scoped>
.sidebar-content {
	display: flex;
	flex-direction: column;
	height: 100%;
	padding: 0;
	margin: 0;
}
.sidebar-user {
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 24px 16px 16px 16px;
	position: relative;
}
.user-info {
	flex: 1;
}
.user-actions {
	position: absolute;
	top: 16px;
	right: 16px;
	display: flex;
	gap: 8px;
}

.theme-toggle-btn,
.logout-btn {
	background: none;
	border: none;
	font-size: 18px;
	cursor: pointer;
	padding: 4px;
	border-radius: 4px;
	transition: background-color 0.2s;
}

.theme-toggle-btn:hover,
.logout-btn:hover {
	background-color: var(--hover-bg);
}
.avatar {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	background: var(--border-color);
	object-fit: cover;
}
.user-name {
	font-weight: bold;
	color: var(--text-primary);
}
.user-status {
	font-size: 12px;
	color: var(--text-muted);
}
.sidebar-search {
	padding: 0 16px 16px 16px;
}
.sidebar-search input {
	width: 100%;
	padding: 8px;
	border-radius: 8px;
	border: 1px solid var(--border-color);
	background: var(--bg-message);
	color: var(--text-primary);
}
.sidebar-search input:focus {
	border-color: var(--selected-border);
	outline: none;
}
.sidebar-actions {
	padding: 0 16px 16px 16px;
}
.sidebar-chats {
	flex: 1 1 0;
	overflow-y: auto;
	padding: 0 0 8px 0;
}
.sidebar-chat {
	display: flex;
	align-items: center;
	gap: 12px;
	padding: 12px 16px;
	cursor: pointer;
	transition: background 0.2s;
}
.sidebar-chat:hover {
	background: var(--hover-bg);
}
.sidebar-chat.selected {
	background: var(--selected-bg) !important;
	border-left: 4px solid var(--selected-border);
	margin-left: -4px;
	box-shadow: 0 0 0 2px #19875433;
}
.chat-avatar {
	width: 36px;
	height: 36px;
	border-radius: 50%;
	object-fit: cover;
	background: var(--avatar-bg);
}
.chat-name-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
}
.delete-chat-btn {
	background: none;
	border: none;
	font-size: 14px;
	opacity: 0;
	transition: opacity 0.2s;
	cursor: pointer;
	padding: 2px 4px;
	border-radius: 4px;
	color: var(--danger-text);
}
.sidebar-chat:hover .delete-chat-btn {
	opacity: 0.7;
}
.delete-chat-btn:hover {
	opacity: 1;
	background: var(--danger-bg);
}
.chat-time {
	font-size: 12px;
	color: var(--text-muted);
	margin-left: 8px;
}
.chat-last-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
}
.chat-unread {
	background: var(--success-bg);
	color: var(--success-text);
	border-radius: 12px;
	font-size: 12px;
	padding: 2px 8px;
	margin-left: 8px;
}
.chat-info {
	flex: 1 1 0;
	min-width: 0;
}
.chat-name {
	font-weight: 500;
	font-size: 15px;
	color: var(--text-primary);
}
.chat-last {
	font-size: 13px;
	color: var(--text-secondary);
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}
.last-message {
	font-size: 13px;
	color: var(--text-secondary);
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	flex: 1;
}
.last-sender {
	font-weight: 500;
	color: var(--text-primary);
}
.no-messages {
	font-size: 13px;
	color: var(--text-muted);
	font-style: italic;
}
.unread-count {
	background: var(--primary-color);
	color: white;
	border-radius: 50%;
	font-size: 11px;
	font-weight: 600;
	padding: 2px 6px;
	margin-left: 8px;
	min-width: 18px;
	height: 18px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.image-indicator {
	color: var(--primary-color);
	font-weight: 500;
}
.sidebar-loading {
	color: var(--text-secondary);
	text-align: center;
	padding: 16px;
}
.sidebar-error {
	color: var(--danger-text);
	text-align: center;
	padding: 16px;
}
.sidebar-empty {
	color: #888;
	text-align: center;
	padding: 16px;
}

/* Edit Profile Styles */
.edit-profile-btn {
	background: none;
	border: none;
	font-size: 14px;
	margin-left: 8px;
	cursor: pointer;
	opacity: 0.7;
	transition: opacity 0.2s;
}

.edit-profile-btn:hover {
	opacity: 1;
}

.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
}

.edit-profile-modal {
	background: #ffffff;
	border-radius: 8px;
	padding: 0;
	width: 90%;
	max-width: 400px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
}

.modal-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 16px 20px;
	border-bottom: 1px solid var(--border-color);
}

.modal-header h3 {
	margin: 0;
	color: var(--text-primary);
	font-size: 18px;
}

.close-btn {
	background: none;
	border: none;
	font-size: 24px;
	cursor: pointer;
	color: var(--text-secondary);
	padding: 0;
	width: 30px;
	height: 30px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.close-btn:hover {
	color: var(--text-primary);
}

.modal-body {
	padding: 20px;
}

.form-group {
	margin-bottom: 16px;
}

.form-group label {
	display: block;
	margin-bottom: 6px;
	color: var(--text-primary);
	font-weight: 500;
}

.form-control {
	width: 100%;
	padding: 8px 12px;
	border: 1px solid var(--border-color);
	border-radius: 4px;
	background: var(--input-bg);
	color: var(--text-primary);
	font-size: 14px;
}

.form-control:focus {
	outline: none;
	border-color: var(--primary-color);
	box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.form-text {
	font-size: 12px;
	margin-top: 4px;
}

.text-muted {
	color: var(--text-secondary);
}

.modal-footer {
	display: flex;
	justify-content: flex-end;
	gap: 8px;
	padding: 16px 20px;
	border-top: 1px solid var(--border-color);
}

.alert {
	padding: 8px 12px;
	border-radius: 4px;
	margin-bottom: 16px;
}

.alert-danger {
	background-color: #f8d7da;
	border: 1px solid #f5c6cb;
	color: #721c24;
}

/* Dark theme adjustments */
.dark-theme .alert-danger {
	background-color: #2d1b1e;
	border-color: #5a2834;
	color: #f8d7da;
}
</style>
