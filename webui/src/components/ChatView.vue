<template>
	<div class="chat-view d-flex flex-column h-100">
		<ChatHeader :chat="chat" />
		<div
			class="messages flex-grow-1 overflow-auto px-3 py-2"
			style="background: var(--bg-message)"
		>
			<MessageBubble
				v-for="msg in messages"
				:key="msg.id"
				:msg="msg"
				:isOwn="msg.own"
				:chat="chat"
				@reaction-changed="$emit('message-sent')"
				@message-deleted="$emit('message-sent')"
			/>
		</div>
		<div
			class="input-area p-3 border-top"
			style="background: var(--bg-surface)"
		>
			<div class="input-group">
				<input
					ref="fileInput"
					type="file"
					accept="image/*"
					style="display: none"
					@change="handleImageUpload"
				/>
				<button
					class="btn btn-outline-secondary"
					@click="$refs.fileInput.click()"
					title="Upload Image"
				>
					📷
				</button>
				<input
					v-model="input"
					type="text"
					class="form-control"
					placeholder="Type a message..."
					@keyup.enter="send"
					style="
						background: var(--bg-message);
						color: var(--text-primary);
						border-color: var(--border-color);
					"
				/>
				<button
					class="btn btn-outline-secondary"
					@click="send"
					:disabled="!input.trim() && !selectedImage"
				>
					Send
				</button>
			</div>

			<!-- Image preview -->
			<div v-if="selectedImage" class="image-preview mt-2">
				<img :src="selectedImage" alt="Preview" class="preview-image" />
				<button
					class="btn btn-sm btn-outline-danger ms-2"
					@click="removeImage"
					title="Remove image"
				>
					✕
				</button>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref } from 'vue';
import ChatHeader from './ChatHeader.vue';
import MessageBubble from './MessageBubble.vue';
import axios from '../services/axios.js';

const props = defineProps({
	chat: Object,
	messages: Array,
});

const emit = defineEmits(['message-sent']);

const input = ref('');
const selectedImage = ref(null);
const fileInput = ref(null);

// Convert file to base64 data URL with compression
function fileToDataURL(file, maxWidth = 800, quality = 0.8) {
	return new Promise((resolve) => {
		const canvas = document.createElement('canvas');
		const ctx = canvas.getContext('2d');
		const img = new Image();

		img.onload = () => {
			// Calculate new dimensions while maintaining aspect ratio
			let { width, height } = img;
			if (width > maxWidth) {
				height = (height * maxWidth) / width;
				width = maxWidth;
			}

			canvas.width = width;
			canvas.height = height;

			// Draw and compress
			ctx.drawImage(img, 0, 0, width, height);
			const compressedDataURL = canvas.toDataURL('image/jpeg', quality);
			resolve(compressedDataURL);
		};

		const reader = new FileReader();
		reader.onload = (e) => {
			img.src = e.target.result;
		};
		reader.readAsDataURL(file);
	});
}

function handleImageUpload(event) {
	const file = event.target.files[0];
	if (file && file.type.startsWith('image/')) {
		// Compress image before setting it
		fileToDataURL(file, 800, 0.8).then((compressedDataURL) => {
			selectedImage.value = compressedDataURL;
		});
	}
}

function removeImage() {
	selectedImage.value = null;
	if (fileInput.value) {
		fileInput.value.value = '';
	}
}

async function send() {
	if ((!input.value.trim() && !selectedImage.value) || !props.chat?.id)
		return;

	const messageContent = input.value.trim();
	const userId = localStorage.getItem('userId');

	try {
		let response;

		// Send message (with or without image)
		response = await axios.post(
			`/users/${userId}/conversations/${props.chat.id}/messages`,
			{
				content: messageContent,
				imageUrl: selectedImage.value || undefined,
			}
		);

		// Clear input and image after successful send
		input.value = '';
		removeImage();

		// Emit event to parent to refresh messages
		emit('message-sent', {
			messageId: response.data,
			content: messageContent,
			imageUrl: selectedImage.value,
			chatId: props.chat.id,
		});
	} catch (error) {
		console.error('Failed to send message:', error);
	}
}
</script>

<style scoped>
.chat-view {
	background: var(--bg-message);
	min-width: 0;
}
.input-area {
	border-top: 1px solid var(--border-color);
}

.image-preview {
	display: flex;
	align-items: center;
}

.preview-image {
	max-width: 200px;
	max-height: 150px;
	border-radius: 8px;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>
