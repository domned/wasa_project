<template>
	<div :class="['message-bubble', isOwn ? 'own' : '']">
		<div class="message-container">
			<div class="bubble p-2 px-3 rounded shadow-sm">
				<div
					v-if="!isOwn && isGroupChat"
					class="sender-name small fw-bold mb-1"
					:style="{ color: senderColor }"
				>
					{{ msg.senderUsername }}
				</div>
				
				<!-- Display image if present -->
				<div v-if="msg.imageUrl" class="message-image mb-2">
					<img 
						:src="msg.imageUrl" 
						alt="Image" 
						class="img-fluid rounded"
						style="max-width: 300px; max-height: 200px; cursor: pointer;"
						@click="openImageModal"
					/>
				</div>
				
				<!-- Display text if present -->
				<span v-if="msg.text">{{ msg.text }}</span>
				<div class="message-footer">
					<span class="time small text-muted">{{ msg.time }}</span>
					<span v-if="isOwn" class="read-status ms-1">
						<span v-if="msg.isRead" class="read-tick" title="Read by recipient">✓✓</span>
						<span v-else class="sent-tick" title="Sent">✓</span>
					</span>
				</div>
			</div>

			<!-- Message actions (delete, forward) -->
			<div v-if="isOwn" class="message-actions">
				<button
					class="btn btn-sm btn-outline-secondary me-1"
					@click="forwardMessage"
					title="Forward message"
				>
					↗️
				</button>
				<button
					class="btn btn-sm btn-outline-danger"
					@click="deleteMessage"
					title="Delete message"
				>
					🗑️
				</button>
			</div>
			<div v-else class="message-actions">
				<button
					class="btn btn-sm btn-outline-secondary"
					@click="forwardMessage"
					title="Forward message"
				>
					↗️
				</button>
			</div>

			<!-- Reactions display -->
			<div v-if="hasReactions" class="reactions mt-1">
				<span
					v-for="(reaction, emoji) in msg.reactions"
					:key="emoji"
					class="reaction-badge"
					:class="{ 'user-reacted': userHasReacted(emoji) }"
					:title="getReactionTooltip(emoji, reaction)"
					@click="toggleReaction(emoji)"
				>
					{{ emoji }} {{ reaction.count }}
				</span>
			</div>

			<!-- Reaction picker -->
			<div class="reaction-actions mt-1">
				<button
					v-if="!showReactionPicker"
					class="btn btn-sm btn-link reaction-btn"
					@click="showReactionPicker = true"
					title="Add reaction"
				>
					😊
				</button>
				<div v-if="showReactionPicker" class="reaction-picker">
					<span
						v-for="emoji in availableEmojis"
						:key="emoji"
						class="emoji-option"
						@click="addReaction(emoji)"
					>
						{{ emoji }}
					</span>
					<button
						class="btn btn-sm btn-link"
						@click="showReactionPicker = false"
					>
						✕
					</button>
				</div>
			</div>
		</div>
	</div>
	
	<!-- Image Modal -->
	<div v-if="showImageModal" class="image-modal-overlay" @click="closeImageModal">
		<div class="image-modal-content" @click.stop>
			<img :src="msg.imageUrl" alt="Full size image" class="modal-image" />
			<button class="close-modal-btn" @click="closeImageModal">✕</button>
		</div>
	</div>
</template>

<script setup>
import { computed, ref } from 'vue';
import axios from '../services/axios.js';

const props = defineProps({
	msg: Object,
	isOwn: Boolean,
	chat: Object,
});

const emit = defineEmits(['reaction-changed', 'message-deleted']);

const showReactionPicker = ref(false);
const showImageModal = ref(false);

function openImageModal() {
	showImageModal.value = true;
}

function closeImageModal() {
	showImageModal.value = false;
}

// Available emoji reactions
const availableEmojis = ['👍', '❤️', '😂', '😮', '😢', '😡'];

// Check if this is a group chat (more than 2 participants)
const isGroupChat = computed(() => {
	return (
		props.chat &&
		props.chat.participants &&
		props.chat.participants.length > 2
	);
});

// Check if message has reactions
const hasReactions = computed(() => {
	return props.msg.reactions && Object.keys(props.msg.reactions).length > 0;
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

// Check if current user has reacted with specific emoji
function userHasReacted(emoji) {
	const currentUsername = getCurrentUsername();
	const reaction = props.msg.reactions[emoji];
	return (
		reaction && reaction.users && reaction.users.includes(currentUsername)
	);
}

// Get current user's username
function getCurrentUsername() {
	// This should match the current logged-in user
	// For now, we'll use a simple approach - in a real app, this would come from auth state
	return localStorage.getItem('currentUsername') || 'Alice'; // Default fallback
}

// Get tooltip text for reaction
function getReactionTooltip(emoji, reaction) {
	if (reaction.count === 1) {
		return `${reaction.users} reacted with ${emoji}`;
	}
	return `${reaction.count} people reacted with ${emoji}: ${reaction.users}`;
}

// Add or remove reaction
async function toggleReaction(emoji) {
	const hasReacted = userHasReacted(emoji);

	if (hasReacted) {
		await removeReaction(emoji);
	} else {
		await addReaction(emoji);
	}
}

// Add reaction to message
async function addReaction(emoji) {
	try {
		const userId = localStorage.getItem('userId');
		await axios.post(
			`/users/${userId}/conversations/${props.chat.id}/messages/${props.msg.id}/reaction`,
			{ emoji: emoji }
		);

		showReactionPicker.value = false;
		emit('reaction-changed');
	} catch (error) {
		console.error('Failed to add reaction:', error);
	}
}

// Remove reaction from message
async function removeReaction(emoji) {
	try {
		const userId = localStorage.getItem('userId');
		await axios.delete(
			`/users/${userId}/conversations/${props.chat.id}/messages/${props.msg.id}/reaction/${emoji}`
		);

		emit('reaction-changed');
	} catch (error) {
		console.error('Failed to remove reaction:', error);
	}
}

// Delete message (only for message owner)
async function deleteMessage() {
	if (!confirm('Are you sure you want to delete this message?')) {
		return;
	}

	try {
		const userId = localStorage.getItem('userId');
		await axios.delete(
			`/users/${userId}/conversations/${props.chat.id}/messages/${props.msg.id}`
		);

		emit('message-deleted', props.msg.id);
	} catch (error) {
		console.error('Failed to delete message:', error);
		alert(
			'Failed to delete message. You can only delete your own messages.'
		);
	}
}

// Forward message to another conversation
async function forwardMessage() {
	try {
		// Get user's conversations to show a selection dialog
		const userId = localStorage.getItem('userId');
		const response = await axios.get(`/users/${userId}/conversations`);
		const conversations = response.data || [];

		// Filter out current conversation
		const otherConversations = conversations.filter(
			(conv) => conv.id !== props.chat.id
		);

		if (otherConversations.length === 0) {
			alert('No other conversations available to forward to.');
			return;
		}

		// Create a selection dialog with conversation names
		let selectionText = 'Select a conversation to forward to:\n\n';
		otherConversations.forEach((conv, index) => {
			const displayName = getConversationDisplayName(conv);
			selectionText += `${index + 1}. ${displayName}\n`;
		});
		selectionText += '\nEnter the number of your choice:';

		const selection = prompt(selectionText);

		if (!selection) {
			return;
		}

		const selectedIndex = parseInt(selection) - 1;

		if (selectedIndex < 0 || selectedIndex >= otherConversations.length) {
			alert('Invalid selection. Please try again.');
			return;
		}

		const targetConversation = otherConversations[selectedIndex];

		await axios.post(
			`/users/${userId}/conversations/${props.chat.id}/messages/${props.msg.id}/forward`,
			{ content: targetConversation.id }
		);

		alert(
			`Message forwarded to "${getConversationDisplayName(
				targetConversation
			)}" successfully!`
		);
	} catch (error) {
		console.error('Failed to forward message:', error);
		alert('Failed to forward message. Please try again.');
	}
}

// Helper function to get conversation display name
function getConversationDisplayName(conversation) {
	const currentUserId = localStorage.getItem('userId');

	// If it's a named group, use the group name
	if (conversation.name && conversation.name.trim()) {
		return conversation.name;
	}

	// For 1:1 chats, show the other participant's name
	if (conversation.participants && conversation.participants.length === 2) {
		const otherParticipant = conversation.participants.find(
			(p) => p.id !== currentUserId
		);
		return otherParticipant ? otherParticipant.username : 'Unknown User';
	}

	// For group chats without a name, show participant names
	if (conversation.participants && conversation.participants.length > 2) {
		const otherParticipants = conversation.participants
			.filter((p) => p.id !== currentUserId)
			.map((p) => p.username)
			.join(', ');
		return `Group: ${otherParticipants}`;
	}

	return 'Unknown Conversation';
}
</script>

<style scoped>
.message-bubble {
	display: flex;
	margin-bottom: 8px;
}
.message-bubble.own {
	justify-content: flex-end;
}
.message-container {
	max-width: 70vw;
}
.bubble {
	background: var(--bg-surface);
	color: var(--text-primary);
	min-width: 40px;
	word-break: break-word;
	border-radius: 18px;
	border: 1px solid var(--border-color);
}
.message-bubble.own .bubble {
	background: var(--bg-own-message);
	color: var(--text-primary);
}

/* Reactions */
.reactions {
	display: flex;
	flex-wrap: wrap;
	gap: 4px;
	margin-left: 8px;
}
.reaction-badge {
	background: var(--bg-message);
	border: 1px solid var(--border-color);
	border-radius: 12px;
	padding: 2px 6px;
	font-size: 12px;
	cursor: pointer;
	transition: all 0.2s;
	color: var(--text-primary);
}
.reaction-badge:hover {
	background: var(--hover-bg);
}
.reaction-badge.user-reacted {
	background: var(--selected-bg);
	border-color: var(--selected-border);
	color: var(--text-primary);
}

/* Reaction actions */
.reaction-actions {
	margin-left: 8px;
}
.reaction-btn {
	padding: 2px 6px;
	font-size: 14px;
	border: none;
	background: transparent;
	opacity: 0.6;
	transition: opacity 0.2s;
}
.reaction-btn:hover {
	opacity: 1;
	background: var(--hover-bg);
}

/* Reaction picker */
.reaction-picker {
	display: flex;
	align-items: center;
	gap: 4px;
	padding: 4px 8px;
	background: var(--bg-surface);
	border: 1px solid var(--border-color);
	border-radius: 8px;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	margin-top: 4px;
}
.emoji-option {
	padding: 4px;
	cursor: pointer;
	border-radius: 4px;
	transition: background 0.2s;
}
.emoji-option:hover {
	background: var(--hover-bg);
}

/* Message actions */
.message-actions {
	margin-top: 4px;
	opacity: 0;
	transition: opacity 0.2s;
	display: flex;
	gap: 4px;
}
.message-container:hover .message-actions {
	opacity: 1;
}
.message-actions .btn {
	padding: 2px 6px;
	font-size: 12px;
	min-width: auto;
}

/* Message footer with time and read status */
.message-footer {
	display: flex;
	align-items: center;
	justify-content: flex-end;
	margin-top: 4px;
	font-size: 12px;
}

/* Read status indicators */
.read-status {
	font-size: 12px;
	margin-left: 4px;
}
.read-tick {
	color: #007bff; /* Blue color for read messages */
	font-weight: bold;
	text-shadow: 0 0 1px rgba(0, 123, 255, 0.5); /* Slight glow effect */
}
.sent-tick {
	color: #6c757d; /* Gray color for sent but unread messages */
}
.unread-tick {
	color: #6c757d; /* Gray color for unread messages */
}

/* Dark theme adjustments */
:global(.dark-theme) .read-tick {
	color: #4dabf7; /* Lighter blue for dark theme */
	text-shadow: 0 0 1px rgba(77, 171, 247, 0.5);
}
:global(.dark-theme) .sent-tick,
:global(.dark-theme) .unread-tick {
	color: #adb5bd; /* Lighter gray for dark theme */
}

/* Image Modal Styles */
.image-modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.8);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 2000;
}

.image-modal-content {
	position: relative;
	max-width: 90%;
	max-height: 90%;
}

.modal-image {
	max-width: 100%;
	max-height: 100%;
	border-radius: 8px;
}

.close-modal-btn {
	position: absolute;
	top: 10px;
	right: 10px;
	background: rgba(0, 0, 0, 0.7);
	color: white;
	border: none;
	border-radius: 50%;
	width: 30px;
	height: 30px;
	font-size: 18px;
	cursor: pointer;
	display: flex;
	align-items: center;
	justify-content: center;
}

.close-modal-btn:hover {
	background: rgba(0, 0, 0, 0.9);
}
</style>
