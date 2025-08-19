<template>
	<div class="sidebar-content">
			<div class="sidebar-user">
				<div class="avatar"></div>
				<div>
					<div class="user-name">{{ username }}</div>
					<div class="user-status">Online</div>
				</div>
			</div>
		<div class="sidebar-search">
			<input type="text" placeholder="Search or start new chat" />
		</div>
		<div class="sidebar-chats">
			<div v-if="loading" class="sidebar-loading">Loading chats...</div>
			<div v-else-if="error" class="sidebar-error">{{ error }}</div>
			<div v-else-if="chats.length === 0" class="sidebar-empty">
				No chats found.
			</div>
			<div
				v-else
				class="sidebar-chat"
				v-for="chat in chats"
				:key="chat.id"
				:class="{ selected: chat.id === selectedChatId }"
				@click="$emit('select-chat', chat.id)"
			>
				<img
					class="chat-avatar"
					:src="getChatAvatar(chat)"
					alt=""
				/>
				<div class="chat-info">
					<div class="chat-name-row">
						<span class="chat-name">
							{{ getChatDisplayName(chat) }}
						</span>
						<!-- Optionally add time if available -->
					</div>
					<div class="chat-last-row">
						<!-- Optionally add last message if available -->
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import axios from '../services/axios.js';
// Returns the avatar for a chat: for 1:1 chats, show the other participant's picture; for groups, show the chat picture
function getChatAvatar(chat) {
	// Group chat: use chat.picture
	if (chat.participants && chat.participants.length > 2) {
		return chat.picture || '';
	}
	// 1:1 chat: use the other participant's picture
	if (chat.participants && chat.participants.length === 2) {
		const other = chat.participants.find((p) => p.id !== props.userId);
		return other && other.picture ? other.picture : chat.picture || '';
	}
	// Fallback
	return chat.picture || '';
}
// Returns the display name for a chat: for 1:1 chats, show the other participant's username; for groups, show the chat name
function getChatDisplayName(chat) {
	// If group chat (more than 2 participants)
	if (chat.participants && chat.participants.length > 2) {
		if (chat.name && chat.name.trim() !== '') {
			return chat.name;
		}
		// Fallback: join all usernames except the current user
		return chat.participants
			.filter((p) => p.id !== props.userId)
			.map((p) => p.username)
			.join(', ');
	}
	// For 1:1 chat, show the other participant's username
	if (chat.participants && chat.participants.length === 2) {
		// Find the participant who is not the current user
		const other = chat.participants.find((p) => p.id !== props.userId);
		return other ? other.username : chat.name || chat.id;
	}
	// Fallback
	return chat.name || chat.id;
}

const emit = defineEmits(['select-chat', 'chats-loaded']);
const props = defineProps({
	userId: {
		type: String,
		required: true,
	},
	username: {
		type: String,
		required: true,
	},
	selectedChatId: {
		type: [String, null],
		required: false,
	},
});

const chats = ref([]);
const loading = ref(false);
const error = ref(null);

async function fetchChats() {
	loading.value = true;
	error.value = null;
	try {
		const res = await axios.get(`/users/${props.userId}/conversations`);
		chats.value = res.data;
			// Emit chats to parent so App.vue can sync selectedChat
			emit('chats-loaded', chats.value);
	} catch (err) {
		error.value = err.message || 'Failed to load chats';
		chats.value = [];
			emit('chats-loaded', []);
	} finally {
		loading.value = false;
	}
}

onMounted(fetchChats);
watch(() => props.userId, fetchChats);
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
}
.avatar {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	background: #b0b0b0;
}
.user-name {
	font-weight: bold;
}
.user-status {
	font-size: 12px;
	color: #888;
}
.sidebar-search {
	padding: 0 16px 16px 16px;
}
.sidebar-search input {
	width: 100%;
	padding: 8px;
	border-radius: 8px;
	border: 1px solid #e0e0e0;
	background: #f5f3ef;
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
	background: #eceae6;
}
.sidebar-chat.selected {
	background: #b6f0e7 !important;
	/* Use a more visible color for highlight */
	border-left: 4px solid #198754;
	/* Bootstrap green */
	margin-left: -4px;
	box-shadow: 0 0 0 2px #19875433;
}
.chat-avatar {
	width: 36px;
	height: 36px;
	border-radius: 50%;
	object-fit: cover;
	background: #d0d0d0;
}
.chat-name-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
}
.chat-time {
	font-size: 12px;
	color: #aaa;
	margin-left: 8px;
}
.chat-last-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
}
.chat-unread {
	background: #4caf50;
	color: #fff;
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
}
.chat-last {
	font-size: 13px;
	color: #888;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}
.sidebar-loading {
	color: #888;
	text-align: center;
	padding: 16px;
}
.sidebar-error {
	color: #c00;
	text-align: center;
	padding: 16px;
}
.sidebar-empty {
	color: #888;
	text-align: center;
	padding: 16px;
}
</style>
