<template>
	<!-- Show login screen if no user is logged in -->
	<Login v-if="!isLoggedIn" @login-success="handleLoginSuccess" />

	<!-- Show main chat interface if user is logged in -->
	<div v-else class="layout-root">
		<aside class="sidebar-root">
			<Sidebar
				ref="sidebarRef"
				:userId="userId"
				:username="username"
				:userPicture="userPicture"
				:selectedChatId="selectedChatId"
				@select-chat="selectChat"
				@chats-loaded="handleChatsLoaded"
				@chat-deleted="handleChatDeleted"
				@logout="handleLogout"
				@username-updated="handleUsernameUpdated"
			/>
		</aside>
		<div class="divider"></div>
		<main class="chat-root">
			<ChatView
				:chat="selectedChat"
				:messages="selectedMessages"
				@message-sent="handleMessageSent"
			/>
		</main>

		<!-- Create Conversation Modal -->
		<CreateConversation @conversation-created="handleConversationCreated" />
	</div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import axios from './services/axios.js';
import Login from './components/Login.vue';
import Sidebar from './components/Sidebar.vue';
import ChatView from './components/ChatView.vue';
import CreateConversation from './components/CreateConversation.vue';

// Check if user is logged in
const isLoggedIn = ref(false);
const userId = ref(null);
const username = ref('');
const userPicture = ref('');
const selectedChatId = ref(null);
const selectedChat = ref(null);
const selectedMessages = ref([]);
const chats = ref([]); // Store all chats from Sidebar
const sidebarRef = ref(null);

// Check login status on app start
function checkLoginStatus() {
	const storedUserId = localStorage.getItem('userId');
	const storedUsername = localStorage.getItem('currentUsername');

	if (storedUserId && storedUsername) {
		userId.value = storedUserId;
		username.value = storedUsername;
		isLoggedIn.value = true;
		fetchUserData();
	}
}

function handleLoginSuccess(userData) {
	userId.value = userData.userId;
	username.value = userData.username;
	isLoggedIn.value = true;
	fetchUserData();
	startChatPolling(); // Start polling after login
}

async function fetchUserData() {
	if (!userId.value) return;
	try {
		const res = await axios.get('/users');
		const user = (res.data || []).find((u) => u.id === userId.value);
		if (user) {
			username.value = user.username;
			userPicture.value = user.picture || '';
			// Update localStorage with current data
			localStorage.setItem('currentUsername', username.value);
		}
	} catch (err) {
		console.error('Failed to fetch user data:', err);
		// If user data fetch fails, force re-login
		handleLogout();
	}
}

function handleLogout() {
	stopChatPolling(); // Stop polling on logout
	localStorage.removeItem('userId');
	localStorage.removeItem('currentUsername');
	userId.value = null;
	username.value = '';
	userPicture.value = '';
	isLoggedIn.value = false;
	selectedChatId.value = null;
	selectedChat.value = null;
	selectedMessages.value = [];
	chats.value = [];
}

function handleUsernameUpdated(newUsername) {
	username.value = newUsername;
	// Update localStorage is already handled in Sidebar component
}

onMounted(() => {
	checkLoginStatus();
	startChatPolling();
});

onUnmounted(() => {
	stopChatPolling();
});

// Polling mechanism to update chats periodically
let pollingInterval = null;
let activePollingInterval = null;

function startChatPolling() {
	// Poll every 15 seconds to refresh sidebar (for new chats and general updates)
	pollingInterval = setInterval(() => {
		if (isLoggedIn.value && sidebarRef.value) {
			sidebarRef.value.refreshChats();
		}
	}, 15000); // 15 seconds
	
	// Poll every 5 seconds for active chat messages (good balance of responsiveness and load)
	activePollingInterval = setInterval(() => {
		if (isLoggedIn.value && selectedChatId.value) {
			selectChat(selectedChatId.value);
		}
	}, 5000); // 5 seconds
}

function stopChatPolling() {
	if (pollingInterval) {
		clearInterval(pollingInterval);
		pollingInterval = null;
	}
	if (activePollingInterval) {
		clearInterval(activePollingInterval);
		activePollingInterval = null;
	}
}

function handleChatsLoaded(loadedChats) {
	chats.value = loadedChats;
}

async function selectChat(chatId) {
	selectedChatId.value = chatId;
	// Find the chat object by id
	selectedChat.value = chats.value.find((c) => c.id === chatId) || null;
	
	// Store previous message count to detect new messages
	const previousMessageCount = selectedMessages.value.length;
	
	try {
		const res = await axios.get(
			`/users/${userId.value}/conversations/${chatId}/messages`
		);
		// Add 'own' property to each message based on current user ID
		const newMessages = res.data.map((msg) => ({
			...msg,
			own: msg.senderId === userId.value,
		}));
		
		// Check if new messages were received (not sent by current user)
		if (newMessages.length > previousMessageCount && previousMessageCount > 0) {
			const latestMessage = newMessages[newMessages.length - 1];
			if (!latestMessage.own && sidebarRef.value) {
				// New message received from someone else, update sidebar
				sidebarRef.value.updateChatWithNewMessage(chatId, {
					id: latestMessage.id,
					senderId: latestMessage.senderId,
					text: latestMessage.text,
					senderUsername: latestMessage.senderUsername
				});
			}
		}
		
		selectedMessages.value = newMessages;
	} catch (err) {
		selectedMessages.value = [];
	}
}

async function handleMessageSent(messageData) {
	// Refresh the messages for the current chat
	if (selectedChatId.value) {
		await selectChat(selectedChatId.value);
	}
	
	// Update the sidebar with the new message to reorder chats
	if (sidebarRef.value && messageData) {
		const newMessage = {
			id: messageData.messageId || Date.now().toString(),
			senderId: userId.value,
			text: messageData.content || '',
			senderUsername: username.value
		};
		
		sidebarRef.value.updateChatWithNewMessage(messageData.chatId, newMessage);
	}
}

async function handleConversationCreated(conversation) {
	// Refresh the sidebar to show the new conversation
	if (sidebarRef.value) {
		await sidebarRef.value.refreshChats();
	}

	// Auto-select the new conversation if it was provided
	if (conversation && conversation.id) {
		// Wait a moment for the sidebar to refresh, then select the new chat
		setTimeout(async () => {
			await selectChat(conversation.id);
		}, 100);
	}
}

function handleChatDeleted(chatId) {
	// If the deleted chat was selected, clear the selection
	if (selectedChatId.value === chatId) {
		selectedChatId.value = null;
		selectedChat.value = null;
		selectedMessages.value = [];
	}
	// Update local chats array
	chats.value = chats.value.filter((chat) => chat.id !== chatId);
}
</script>

<style>
/* Light theme (default) */
:root {
	--bg-primary: #edeae3;
	--bg-secondary: #f8f6f3;
	--bg-surface: #ffffff;
	--bg-chat: rgb(168, 155, 129);
	--bg-message: #f5f3ef;
	--bg-own-message: #e7e3d8;
	--bg-overlay: rgba(255, 255, 255, 0.9);
	--text-primary: #333;
	--text-secondary: #6b5e3a;
	--text-muted: #888;
	--border-color: #d1d1d1;
	--border-subtle: #e0d8c8;
	--hover-bg: #eceae6;
	--selected-bg: #b6f0e7;
	--selected-border: #198754;
	--avatar-bg: #007bff;
	--avatar-text: #ffffff;
	--danger-bg: rgba(220, 53, 69, 0.1);
	--danger-text: #dc3545;
	--danger-hover: #c82333;
	--success-bg: #4caf50;
	--success-text: #ffffff;
}

/* Dark theme */
.dark-theme {
	--bg-primary: #1a1a1a;
	--bg-secondary: #2d2d2d;
	--bg-surface: #3a3a3a;
	--bg-chat: #333333;
	--bg-message: #404040;
	--bg-own-message: #4a4a4a;
	--bg-overlay: rgba(45, 45, 45, 0.9);
	--text-primary: #ffffff;
	--text-secondary: #cccccc;
	--text-muted: #999999;
	--border-color: #555555;
	--border-subtle: #444444;
	--hover-bg: #3a3a3a;
	--selected-bg: #2a4a47;
	--selected-border: #4caf50;
	--avatar-bg: #007bff;
	--avatar-text: #ffffff;
	--danger-bg: rgba(220, 53, 69, 0.2);
	--danger-text: #ff6b6b;
	--danger-hover: rgba(220, 53, 69, 0.3);
	--success-bg: #4caf50;
	--success-text: #ffffff;
}

.layout-root {
	display: flex;
	height: 100vh;
	width: 100vw;
	background: var(--bg-primary);
	overflow: hidden;
	color: var(--text-primary);
}
.sidebar-root {
	width: 220px;
	min-width: 220px;
	max-width: 220px;
	height: 100vh;
	background: var(--bg-secondary);
	display: flex;
	flex-direction: column;
	box-shadow: none;
	border: none;
	padding: 0;
	margin: 0;
}
.divider {
	width: 1px;
	background: var(--border-color);
	height: 100vh;
	align-self: stretch;
}
.chat-root {
	flex: 1 1 0;
	min-width: 0;
	height: 100vh;
	display: flex;
	flex-direction: column;
	background: var(--bg-chat);
	padding: 0;
	margin: 0;
}

.user-picker {
	padding: 12px 8px 0 8px;
	background: var(--bg-secondary);
	border-bottom: 1px solid var(--border-subtle);
	font-size: 15px;
	color: var(--text-secondary);
	display: flex;
	align-items: center;
	gap: 8px;
}
.user-picker select {
	font-size: 15px;
	padding: 2px 8px;
	border-radius: 4px;
	border: 1px solid var(--border-color);
	background: var(--bg-primary);
	color: var(--text-primary);
}
.db-footer {
	width: 100vw;
	background: var(--bg-secondary);
	color: var(--text-secondary);
	text-align: center;
	font-size: 14px;
	padding: 6px 0;
	border-top: 1px solid var(--border-color);
	position: fixed;
	bottom: 0;
	left: 0;
	z-index: 100;
}
</style>
console.log("App loaded");
