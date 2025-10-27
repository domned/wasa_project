<template>
	<div class="modal fade" id="createConversationModal" tabindex="-1">
		<div class="modal-dialog">
			<div class="modal-content">
				<!-- Loading overlay -->
				<div v-if="isCreating" class="loading-overlay">
					<div class="spinner-border text-primary" role="status">
						<span class="visually-hidden">Creating...</span>
					</div>
					<div class="mt-2">Creating conversation...</div>
				</div>

				<div class="modal-header">
					<h5 class="modal-title">Create New Conversation</h5>
					<button
						type="button"
						class="btn-close"
						@click="resetAndClose"
						:disabled="isCreating"
					></button>
				</div>
				<div class="modal-body">
					<div class="alert alert-info" role="alert">
						<small
							>💡 Add one person for a private chat, or multiple
							people to create a group!</small
						>
					</div>
					<div
						v-if="availableUsers.length === 0"
						class="alert alert-danger"
						role="alert"
					>
						<strong>⚠️ No users available!</strong>
						<p class="mb-0 mt-2">
							Error in DB creation. Please run
							<code>docker-compose down -v</code> and rebuild the
							container.
						</p>
					</div>
					<form @submit.prevent="createConversation">
						<div class="mb-3">
							<label for="conversationName" class="form-label">
								Conversation Name
								<span
									v-if="participants.length > 1"
									class="text-danger"
									>*</span
								>
								<span v-else class="text-muted"
									>(Optional for 1-on-1 chats)</span
								>
							</label>
							<input
								v-model="conversationName"
								type="text"
								class="form-control"
								id="conversationName"
								:placeholder="
									participants.length > 1
										? 'Enter group name'
										: 'Enter conversation name'
								"
								:required="participants.length > 1"
							/>
							<small
								v-if="participants.length > 1"
								class="form-text text-muted"
							>
								Group name is required for conversations with
								multiple participants.
							</small>
						</div>
						<div class="mb-3">
							<label for="participants" class="form-label"
								>Add Participants</label
							>
							<select
								v-model="selectedUser"
								class="form-select"
								@change="addParticipant"
							>
								<option value="">
									Select a user to add...
								</option>
								<option
									v-for="user in availableUsers"
									:key="user.id"
									:value="user.id"
								>
									{{ user.username }}
								</option>
							</select>
						</div>
						<div class="mb-3" v-if="participants.length > 0">
							<label class="form-label"
								>Selected Participants:</label
							>
							<div class="participant-list">
								<span
									v-for="participant in participants"
									:key="participant.id"
									class="badge bg-primary me-2 mb-2"
								>
									{{ participant.username }}
									<button
										type="button"
										class="btn-close btn-close-white ms-2"
										@click="
											removeParticipant(participant.id)
										"
									></button>
								</span>
							</div>
						</div>
						<div class="modal-footer">
							<button
								type="button"
								class="btn btn-secondary"
								@click="resetAndClose"
								:disabled="isCreating"
							>
								Cancel
							</button>
							<button
								type="submit"
								class="btn btn-primary"
								:disabled="
									participants.length === 0 ||
									isCreating ||
									(participants.length > 1 &&
										!conversationName.trim())
								"
							>
								<span v-if="isCreating">Creating...</span>
								<span v-else-if="participants.length > 1"
									>Create Group</span
								>
								<span v-else-if="participants.length === 1"
									>Start Chat</span
								>
								<span v-else>Select Participants</span>
							</button>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../services/api.js';

const emit = defineEmits(['conversation-created']);

const conversationName = ref('');
const selectedUser = ref('');
const participants = ref([]);
const availableUsers = ref([]);
const isCreating = ref(false);

const loadUsers = async () => {
	try {
		const users = await api.users.listAll();
		const currentUserId = localStorage.getItem('userId');
		// Filter out current user from available users
		availableUsers.value = users.filter(
			(user) => user.id !== currentUserId
		);
	} catch (error) {
		console.error('Failed to load users:', error);
	}
};

const addParticipant = () => {
	if (!selectedUser.value) return;

	const user = availableUsers.value.find((u) => u.id === selectedUser.value);
	if (user && !participants.value.find((p) => p.id === user.id)) {
		participants.value.push(user);
	}
	selectedUser.value = '';
};

const removeParticipant = (userId) => {
	participants.value = participants.value.filter((p) => p.id !== userId);
};

const createConversation = async () => {
	if (participants.value.length === 0) return;

	// Validate group name for multi-participant conversations
	if (participants.value.length > 1 && !conversationName.value.trim()) {
		alert(
			'Please enter a group name for conversations with multiple participants.'
		);
		return;
	}

	isCreating.value = true;

	try {
		const userId = localStorage.getItem('userId');
		const participantIds = participants.value.map((p) => p.id);

		// Include the current user in the participants list
		if (!participantIds.includes(userId)) {
			participantIds.push(userId);
		}

		const requestData = {
			participants: participantIds,
			name: conversationName.value || undefined,
		};

		console.log('Creating conversation with data:', requestData);

		const conversation = await api.conversations.create(
			userId,
			participantIds,
			conversationName.value || undefined
		);

		console.log('Conversation created successfully:', conversation);

		// Close modal properly
		const modal = document.getElementById('createConversationModal');
		const existingModal = bootstrap.Modal.getInstance(modal);
		if (existingModal) {
			existingModal.hide();
		} else {
			// Fallback: create new instance and hide
			const bsModal = new bootstrap.Modal(modal);
			bsModal.hide();
		}

		// Additional fallback with timeout
		setTimeout(() => {
			if (modal.style.display !== 'none') {
				modal.style.display = 'none';
				document.body.classList.remove('modal-open');
				const backdrop = document.querySelector('.modal-backdrop');
				if (backdrop) {
					backdrop.remove();
				}
			}
		}, 500);

		// Reset form
		conversationName.value = '';
		participants.value = [];
		selectedUser.value = '';

		// Emit event to parent
		emit('conversation-created', conversation);
	} catch (error) {
		console.error('Failed to create conversation:', error);
		const resp = error.response;
		const msg = resp?.data?.message || resp?.data || '';
		let friendly = 'Failed to create conversation.';
		if (resp) {
			if (resp.status === 400) {
				friendly =
					msg ||
					'Invalid request. Please check participants and name.';
			} else if (resp.status === 404) {
				friendly = msg || 'One or more participants were not found.';
			} else {
				friendly =
					msg || `Server error (${resp.status}). Please try again.`;
			}
		} else {
			friendly =
				'Network error. Please check your connection and try again.';
		}
		alert(friendly);

		// Make sure to reset the creating state even on error
		isCreating.value = false;
	} finally {
		// Ensure the loading state is always reset
		isCreating.value = false;
	}
};

// Reset form and close modal
const resetAndClose = () => {
	// Reset all form data
	conversationName.value = '';
	participants.value = [];
	selectedUser.value = '';
	isCreating.value = false;

	// Force close the modal
	const modal = document.getElementById('createConversationModal');
	const existingModal = bootstrap.Modal.getInstance(modal);
	if (existingModal) {
		existingModal.hide();
	} else {
		// Fallback: use data-bs-dismiss
		modal.style.display = 'none';
		document.body.classList.remove('modal-open');
		const backdrop = document.querySelector('.modal-backdrop');
		if (backdrop) {
			backdrop.remove();
		}
	}
};

onMounted(() => {
	loadUsers();
});
</script>

<style scoped>
.participant-list {
	min-height: 40px;
}

.badge {
	display: inline-flex;
	align-items: center;
}

.btn-close-white {
	filter: invert(1);
}

.loading-overlay {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: var(--bg-overlay);
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	z-index: 1050;
	border-radius: 0.375rem;
	color: var(--text-primary);
}
</style>
