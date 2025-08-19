<template>
	<div :class="['message-bubble', isOwn ? 'own' : '']">
		<div class="bubble p-2 px-3 rounded shadow-sm">
			<div v-if="!isOwn && isGroupChat" class="sender-name small fw-bold mb-1" :style="{ color: senderColor }">
				{{ msg.senderUsername }}
			</div>
			<span>{{ msg.text }}</span>
			<span class="time small text-muted ms-2">{{ msg.time }}</span>
		</div>
	</div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
	msg: Object,
	isOwn: Boolean,
	chat: Object,
});

// Check if this is a group chat (more than 2 participants)
const isGroupChat = computed(() => {
	return props.chat && props.chat.participants && props.chat.participants.length > 2;
});

// Generate a consistent color for each sender based on their username
const senderColor = computed(() => {
	if (!props.msg.senderUsername) return '#666';
	
	// Simple hash function to generate a color from username
	let hash = 0;
	for (let i = 0; i < props.msg.senderUsername.length; i++) {
		hash = props.msg.senderUsername.charCodeAt(i) + ((hash << 5) - hash);
	}
	
	// Convert hash to HSL color with good saturation and lightness
	const hue = Math.abs(hash) % 360;
	return `hsl(${hue}, 65%, 45%)`;
});
</script>

<style scoped>
.message-bubble {
	display: flex;
	margin-bottom: 8px;
}
.message-bubble.own {
	justify-content: flex-end;
}
.bubble {
	background: #fff;
	color: #333;
	max-width: 70vw;
	min-width: 40px;
	word-break: break-word;
	border-radius: 18px;
	border: 1px solid #e0e0e0;
}
.message-bubble.own .bubble {
	background: #e7e3d8;
	color: #222;
}
</style>
