<template>
	<div class="layout-root">
		<aside class="sidebar-root">
			<Sidebar
				:userId="userId"
				:username="username"
				:selectedChatId="selectedChatId"
				@select-chat="selectChat"
				@chats-loaded="handleChatsLoaded"
			/>
		</aside>
		<div class="divider"></div>
		<main class="chat-root">
			<ChatView :chat="selectedChat" :messages="selectedMessages" />
		</main>
	</div>
</template>

<script setup>

import { ref, onMounted } from 'vue';
import axios from './services/axios.js';
import Sidebar from './components/Sidebar.vue';
import ChatView from './components/ChatView.vue';


// Always set a default userId in localStorage if not present (simulate always-logged-in user)
const DEFAULT_USER_ID = 'f2555a8a-2e66-4326-9588-20e7e298d615'; // Set to your login user (Alice)
if (!localStorage.getItem('userId')) {
	localStorage.setItem('userId', DEFAULT_USER_ID);
}
const userId = ref(localStorage.getItem('userId'));
const username = ref('');
	const selectedChatId = ref(null);
	const selectedChat = ref(null);
	const selectedMessages = ref([]);
	const chats = ref([]); // Store all chats from Sidebar

async function fetchUsername() {
	if (!userId.value) return;
	try {
		const res = await axios.get('/users');
		const user = (res.data || []).find(u => u.id === userId.value);
		username.value = user ? user.username : '';
	} catch (err) {
		username.value = '';
	}
}

onMounted(() => {
	fetchUsername();
});


	function handleChatsLoaded(loadedChats) {
		chats.value = loadedChats;
	}

	async function selectChat(chatId) {
		selectedChatId.value = chatId;
		// Find the chat object by id
		selectedChat.value = chats.value.find(c => c.id === chatId) || null;
		try {
			   const res = await axios.get(`/users/${userId.value}/conversations/${chatId}/messages`);
			   // Add 'own' property to each message based on current user ID
			   selectedMessages.value = res.data.map(msg => ({
				   ...msg,
				   own: msg.senderId === userId.value
			   }));
		} catch (err) {
			selectedMessages.value = [];
		}
	}


</script>

<style>
.layout-root {
	display: flex;
	height: 100vh;
	width: 100vw;
	background: #edeae3;
	overflow: hidden;
}
.sidebar-root {
	width: 220px;
	min-width: 220px;
	max-width: 220px;
	height: 100vh;
	background: #f8f6f3;
	display: flex;
	flex-direction: column;
	box-shadow: none;
	border: none;
	padding: 0;
	margin: 0;
}
.divider {
	width: 1px;
	background: #d1d1d1;
	height: 100vh;
	align-self: stretch;
}
.chat-root {
	flex: 1 1 0;
	min-width: 0;
	height: 100vh;
	display: flex;
	flex-direction: column;
	background: rgb(168, 155, 129);
	/* stronger beige */
	padding: 0;
	margin: 0;
}

.user-picker {
	padding: 12px 8px 0 8px;
	background: #f8f6f3;
	border-bottom: 1px solid #e0d8c8;
	font-size: 15px;
	color: #6b5e3a;
	display: flex;
	align-items: center;
	gap: 8px;
}
.user-picker select {
	font-size: 15px;
	padding: 2px 8px;
	border-radius: 4px;
	border: 1px solid #d1d1d1;
}
.db-footer {
	width: 100vw;
	background: #f5e9c8;
	color: #6b5e3a;
	text-align: center;
	font-size: 14px;
	padding: 6px 0;
	border-top: 1px solid #d1d1d1;
	position: fixed;
	bottom: 0;
	left: 0;
	z-index: 100;
}
</style>
console.log("App loaded");
