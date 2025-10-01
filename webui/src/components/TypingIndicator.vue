<template>
	<div v-if="typingUsers.length > 0" class="typing-indicator">
		<div class="typing-text">
			<span v-if="typingUsers.length === 1">
				{{ typingUsers[0] }} is typing
			</span>
			<span v-else-if="typingUsers.length === 2">
				{{ typingUsers[0] }} and {{ typingUsers[1] }} are typing
			</span>
			<span v-else>
				{{ typingUsers.slice(0, -1).join(', ') }} and
				{{ typingUsers[typingUsers.length - 1] }} are typing
			</span>
			<div class="typing-dots">
				<span></span>
				<span></span>
				<span></span>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import webSocketService from '../services/websocket.js';

const props = defineProps({
	conversationId: String,
	currentUserId: String,
});

const typingUsers = ref([]);
const typingTimeouts = new Map();

onMounted(() => {
	webSocketService.on('typingStart', handleTypingStart);
	webSocketService.on('typingStop', handleTypingStop);
});

onUnmounted(() => {
	webSocketService.off('typingStart', handleTypingStart);
	webSocketService.off('typingStop', handleTypingStop);
	// Clear all timeouts
	typingTimeouts.forEach((timeout) => clearTimeout(timeout));
	typingTimeouts.clear();
});

function handleTypingStart(data) {
	if (
		data.conversation_id !== props.conversationId ||
		data.user_id === props.currentUserId
	) {
		return;
	}

	const username = data.username;
	if (!typingUsers.value.includes(username)) {
		typingUsers.value.push(username);
	}

	// Clear existing timeout for this user
	if (typingTimeouts.has(username)) {
		clearTimeout(typingTimeouts.get(username));
	}

	// Set timeout to remove user from typing if no stop event received
	const timeout = setTimeout(() => {
		handleTypingStop({ username });
	}, 3000); // Remove after 3 seconds of inactivity

	typingTimeouts.set(username, timeout);
}

function handleTypingStop(data) {
	if (data.conversation_id && data.conversation_id !== props.conversationId) {
		return;
	}

	const username = data.username;
	const index = typingUsers.value.indexOf(username);
	if (index > -1) {
		typingUsers.value.splice(index, 1);
	}

	// Clear timeout for this user
	if (typingTimeouts.has(username)) {
		clearTimeout(typingTimeouts.get(username));
		typingTimeouts.delete(username);
	}
}
</script>

<style scoped>
.typing-indicator {
	padding: 0.5rem 1rem;
	background: var(--bg-surface);
	border-bottom: 1px solid var(--border-color);
	font-size: 0.875rem;
	color: var(--text-muted);
}

.typing-text {
	display: flex;
	align-items: center;
	gap: 0.5rem;
}

.typing-dots {
	display: flex;
	gap: 2px;
}

.typing-dots span {
	width: 4px;
	height: 4px;
	background: var(--text-muted);
	border-radius: 50%;
	animation: typing-bounce 1.4s infinite ease-in-out;
}

.typing-dots span:nth-child(1) {
	animation-delay: -0.32s;
}

.typing-dots span:nth-child(2) {
	animation-delay: -0.16s;
}

@keyframes typing-bounce {
	0%,
	80%,
	100% {
		transform: scale(0.8);
		opacity: 0.5;
	}
	40% {
		transform: scale(1);
		opacity: 1;
	}
}
</style>
