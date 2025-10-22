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
						style="
							max-width: 300px;
							max-height: 200px;
							cursor: pointer;
						"
						@click="openImageModal"
					/>
				</div>

				<!-- Display text if present -->
				<span v-if="msg.text">{{ msg.text }}</span>
				<div class="message-footer">
					<span class="time small text-muted">{{ msg.time }}</span>
					<span v-if="isOwn" class="read-status ms-1">
						<span
							v-if="msg.isRead"
							class="read-tick"
							title="Read by recipient"
							>✓✓</span
						>
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

			<!-- Comments section -->
			<div class="comments-section mt-2">
				<!-- Comments toggle button -->
				<button
					class="btn btn-sm btn-link comments-toggle"
					@click="toggleComments"
					:title="showComments ? 'Hide comments' : 'Show comments'"
				>
					💬 {{ commentCount }}
					{{ commentCount === 1 ? 'comment' : 'comments' }}
				</button>

				<!-- Comments list -->
				<div v-if="showComments" class="comments-list mt-2">
					<div
						v-for="comment in comments"
						:key="comment.id"
						class="comment-item"
					>
						<div class="comment-header">
							<span class="comment-author">{{
								comment.author.username
							}}</span>
							<span class="comment-time">{{
								formatTime(comment.timestamp)
							}}</span>
							<button
								v-if="canDeleteComment(comment)"
								class="btn btn-sm btn-link comment-delete"
								@click="deleteComment(comment.id)"
								title="Delete comment"
							>
								✕
							</button>
						</div>
						<div class="comment-text">{{ comment.text }}</div>
					</div>

					<!-- Add comment form -->
					<div class="add-comment-form mt-2">
						<div class="comment-input-group">
							<input
								v-model="newCommentText"
								type="text"
								class="form-control form-control-sm"
								placeholder="Add a comment..."
								@keyup.enter="addComment"
								maxlength="500"
							/>
							<button
								class="btn btn-sm btn-primary"
								@click="addComment"
								:disabled="
									!newCommentText.trim() || isAddingComment
								"
							>
								{{ isAddingComment ? 'Adding...' : 'Send' }}
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Image Modal -->
	<div
		v-if="showImageModal"
		class="image-modal-overlay"
		@click="closeImageModal"
	>
		<div class="image-modal-content" @click.stop>
			<img
				:src="msg.imageUrl"
				alt="Full size image"
				class="modal-image"
			/>
			<button class="close-modal-btn" @click="closeImageModal">✕</button>
		</div>
	</div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import { apiService } from '../services/api.js';

const props = defineProps({
	msg: Object,
	isOwn: Boolean,
	chat: Object,
});

const emit = defineEmits([
	'reaction-changed',
	'message-deleted',
	'comment-added',
	'comment-deleted',
]);

const showReactionPicker = ref(false);
const showImageModal = ref(false);
const showComments = ref(false);
const comments = ref([]);
const newCommentText = ref('');
const isAddingComment = ref(false);
const currentUserId = localStorage.getItem('userId');

// Initialize comments from message data when component mounts or message changes
function initializeComments() {
	if (props.msg.comments && Array.isArray(props.msg.comments)) {
		comments.value = props.msg.comments;
	} else {
		comments.value = [];
	}
}

// Watch for message prop changes and reinitialize comments
watch(() => props.msg, initializeComments, { immediate: true, deep: true });

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
		await apiService.reactions.toggle(
			userId,
			props.chat.id,
			props.msg.id,
			emoji
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
		await apiService.reactions.remove(
			userId,
			props.chat.id,
			props.msg.id,
			emoji
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
		await apiService.messages.delete(userId, props.chat.id, props.msg.id);

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
		const conversations =
			await apiService.conversations.getUserConversations(userId);

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

		await apiService.messages.forward(
			userId,
			props.chat.id,
			props.msg.id,
			targetConversation.id
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

	// For 1:1 chats, always show the other participant's name (ignore backend-generated names)
	if (conversation.participants && conversation.participants.length === 2) {
		const otherParticipant = conversation.participants.find(
			(p) => p.id !== currentUserId
		);
		return otherParticipant ? otherParticipant.username : 'Unknown User';
	}

	// For group chats (3+ participants), use the group name if provided
	if (conversation.participants && conversation.participants.length > 2) {
		if (conversation.name && conversation.name.trim()) {
			return conversation.name;
		}
		// Otherwise show participant names
		const otherParticipants = conversation.participants
			.filter((p) => p.id !== currentUserId)
			.map((p) => p.username)
			.join(', ');
		return `Group: ${otherParticipants}`;
	}

	return 'Unknown Conversation';
}

// Comment-related computed properties
const commentCount = computed(() => {
	return comments.value.length;
});

// Comment-related methods
function toggleComments() {
	showComments.value = !showComments.value;
	// Comments are already loaded from props.msg.comments via watcher
}

async function loadComments() {
	// Comments are loaded as part of message data from the backend
	// The message object should include a comments array with comment objects
	if (props.msg.comments && Array.isArray(props.msg.comments)) {
		comments.value = props.msg.comments;
	} else {
		// Initialize empty comments array if not present in message data
		comments.value = [];
	}
}

async function addComment() {
	if (!newCommentText.value.trim()) return;

	isAddingComment.value = true;

	try {
		const response = await apiService.comments.add(
			currentUserId,
			props.chat.id,
			props.msg.id,
			newCommentText.value.trim()
		);

		comments.value.push(response);
		newCommentText.value = '';

		// Emit event to notify parent component of new comment
		emit('comment-added', { messageId: props.msg.id, comment: response });
	} catch (error) {
		console.error('Failed to add comment:', error);
	} finally {
		isAddingComment.value = false;
	}
}

async function deleteComment(commentId) {
	if (!confirm('Delete this comment?')) return;

	try {
		await apiService.comments.delete(
			currentUserId,
			props.chat.id,
			props.msg.id,
			commentId
		);

		comments.value = comments.value.filter((c) => c.id !== commentId);

		// Emit event to notify parent component of comment deletion
		emit('comment-deleted', { messageId: props.msg.id, commentId });
	} catch (error) {
		console.error('Failed to delete comment:', error);
	}
}

function canDeleteComment(comment) {
	return comment.author.id === currentUserId;
}

function formatTime(timestamp) {
	if (!timestamp) return '';
	const date = new Date(timestamp);
	return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
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

/* Comments Section */
.comments-section {
	margin-left: 8px;
	border-left: 2px solid var(--border-color);
	padding-left: 8px;
}

.comments-toggle {
	color: var(--text-muted);
	text-decoration: none;
	font-size: 0.75rem;
	padding: 0.25rem 0;
}

.comments-toggle:hover {
	color: var(--text-primary);
	background: transparent;
}

.comments-list {
	max-height: 200px;
	overflow-y: auto;
}

.comment-item {
	background: var(--bg-message);
	border: 1px solid var(--border-color);
	border-radius: 8px;
	padding: 0.5rem;
	margin-bottom: 0.5rem;
}

.comment-header {
	display: flex;
	align-items: center;
	gap: 0.5rem;
	margin-bottom: 0.25rem;
}

.comment-author {
	font-weight: 600;
	font-size: 0.75rem;
	color: var(--text-primary);
}

.comment-time {
	font-size: 0.625rem;
	color: var(--text-muted);
}

.comment-delete {
	margin-left: auto;
	padding: 0;
	color: var(--text-muted);
	font-size: 0.75rem;
}

.comment-delete:hover {
	color: var(--danger);
	background: transparent;
}

.comment-text {
	font-size: 0.8rem;
	color: var(--text-primary);
	word-break: break-word;
}

.add-comment-form {
	border-top: 1px solid var(--border-color);
	padding-top: 0.5rem;
}

.comment-input-group {
	display: flex;
	gap: 0.5rem;
}

.comment-input-group .form-control {
	flex: 1;
	font-size: 0.75rem;
}

.comment-input-group .btn {
	font-size: 0.75rem;
	padding: 0.25rem 0.75rem;
}
</style>
