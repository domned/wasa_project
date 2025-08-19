<template>
	<div class="chat-view d-flex flex-column h-100">
		<ChatHeader :chat="chat" />
		<div
			class="messages flex-grow-1 overflow-auto px-3 py-2"
			style="background: #f5f3ef"
		>
			<MessageBubble
				v-for="msg in messages"
				:key="msg.id"
				:msg="msg"
				:isOwn="msg.own"
				:chat="chat"
			/>
		</div>
		<div class="input-area p-3 bg-white border-top">
			<div class="input-group">
				<input
					v-model="input"
					type="text"
					class="form-control"
					placeholder="Type a message..."
				/>
				<button class="btn btn-outline-secondary" @click="send">
					Send
				</button>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref } from 'vue';
import ChatHeader from './ChatHeader.vue';
import MessageBubble from './MessageBubble.vue';

const props = defineProps({
	chat: Object,
	messages: Array
});

const input = ref('');
function send() {
	// Sending logic can be implemented here
	input.value = '';
}
</script>

<style scoped>
.chat-view {
	background: #f5f3ef;
	min-width: 0;
}
.input-area {
	border-top: 1px solid #e0e0e0;
}
</style>
