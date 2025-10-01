<template>
	<div class="contacts-container">
		<div class="contacts-header">
			<h3>Contacts</h3>
			<button
				class="btn btn-primary btn-sm"
				@click="showAddContact = true"
			>
				+ Add Contact
			</button>
		</div>

		<!-- Loading state -->
		<div v-if="loading" class="contacts-loading">Loading contacts...</div>

		<!-- Error state -->
		<div v-if="error" class="alert alert-danger">
			{{ error }}
		</div>

		<!-- Empty state -->
		<div
			v-if="!loading && !error && contacts.length === 0"
			class="contacts-empty"
		>
			<p>No contacts found.</p>
			<p>Add contacts to start conversations easily!</p>
		</div>

		<!-- Contacts list -->
		<div
			v-if="!loading && !error && contacts.length > 0"
			class="contacts-list"
		>
			<div
				v-for="contact in contacts"
				:key="contact.id"
				class="contact-item"
			>
				<div class="contact-avatar">
					{{ contact.username.charAt(0).toUpperCase() }}
				</div>
				<div class="contact-info">
					<div class="contact-name">{{ contact.username }}</div>
					<div class="contact-id">{{ contact.contactUserId }}</div>
				</div>
				<div class="contact-actions">
					<button
						class="btn btn-sm btn-outline-primary"
						@click="startConversation(contact)"
						title="Start conversation"
					>
						💬
					</button>
					<button
						class="btn btn-sm btn-outline-danger"
						@click="removeContact(contact)"
						title="Remove contact"
					>
						🗑️
					</button>
				</div>
			</div>
		</div>

		<!-- Add Contact Modal -->
		<div
			v-if="showAddContact"
			class="modal-overlay"
			@click="showAddContact = false"
		>
			<div class="add-contact-modal" @click.stop>
				<div class="modal-header">
					<h3>Add New Contact</h3>
					<button class="close-btn" @click="cancelAddContact">
						×
					</button>
				</div>
				<div class="modal-body">
					<div class="form-group">
						<label for="contactSearch">Search Users</label>
						<input
							id="contactSearch"
							type="text"
							v-model="searchQuery"
							placeholder="Search by username..."
							class="form-control"
							@input="searchUsers"
						/>
					</div>

					<!-- Search Results -->
					<div v-if="searchResults.length > 0" class="search-results">
						<h4>Search Results</h4>
						<div
							v-for="user in searchResults"
							:key="user.id"
							class="search-result-item"
							@click="selectUserToAdd(user)"
						>
							<div class="user-avatar">
								{{ user.username.charAt(0).toUpperCase() }}
							</div>
							<div class="user-info">
								<div class="user-name">{{ user.username }}</div>
								<div class="user-id">{{ user.id }}</div>
							</div>
						</div>
					</div>

					<!-- Selected user -->
					<div v-if="selectedUser" class="selected-user">
						<h4>Add to Contacts</h4>
						<div class="user-card">
							<div class="user-avatar">
								{{
									selectedUser.username
										.charAt(0)
										.toUpperCase()
								}}
							</div>
							<div class="user-info">
								<div class="user-name">
									{{ selectedUser.username }}
								</div>
								<div class="user-id">{{ selectedUser.id }}</div>
							</div>
						</div>
					</div>

					<div v-if="addError" class="alert alert-danger">
						{{ addError }}
					</div>
				</div>
				<div class="modal-footer">
					<button class="btn btn-secondary" @click="cancelAddContact">
						Cancel
					</button>
					<button
						class="btn btn-primary"
						@click="addContact"
						:disabled="isAdding || !selectedUser"
					>
						{{ isAdding ? 'Adding...' : 'Add Contact' }}
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../services/api.js';

const props = defineProps({
	userId: {
		type: String,
		required: true,
	},
});

const emit = defineEmits(['conversation-created']);

// State
const contacts = ref([]);
const loading = ref(false);
const error = ref('');

// Add contact modal state
const showAddContact = ref(false);
const searchQuery = ref('');
const searchResults = ref([]);
const selectedUser = ref(null);
const isAdding = ref(false);
const addError = ref('');

// Load contacts on mount
onMounted(() => {
	loadContacts();
});

// Load user's contacts
async function loadContacts() {
	loading.value = true;
	error.value = '';

	try {
		contacts.value = await api.contacts.list(props.userId);
	} catch (err) {
		error.value = 'Failed to load contacts';
		console.error('Failed to load contacts:', err);
	} finally {
		loading.value = false;
	}
}

// Search users for adding as contacts
async function searchUsers() {
	if (searchQuery.value.trim().length < 2) {
		searchResults.value = [];
		return;
	}

	try {
		const allUsers = await api.users.listAll();
		const currentUserId = props.userId;

		// Filter users by search query and exclude current user and existing contacts
		const existingContactIds = contacts.value.map((c) => c.contactUserId);

		searchResults.value = allUsers
			.filter(
				(user) =>
					user.id !== currentUserId &&
					!existingContactIds.includes(user.id) &&
					user.username
						.toLowerCase()
						.includes(searchQuery.value.toLowerCase())
			)
			.slice(0, 10); // Limit to 10 results
	} catch (err) {
		console.error('Failed to search users:', err);
	}
}

// Select a user to add as contact
function selectUserToAdd(user) {
	selectedUser.value = user;
	searchResults.value = [];
	searchQuery.value = user.username;
}

// Add selected user as contact
async function addContact() {
	if (!selectedUser.value) return;

	isAdding.value = true;
	addError.value = '';

	try {
		const newContact = await api.contacts.add(
			props.userId,
			selectedUser.value.id
		);

		// Add to local contacts list
		contacts.value.push(newContact);

		// Close modal
		cancelAddContact();
	} catch (err) {
		addError.value =
			'Failed to add contact. User may already be in your contacts.';
		console.error('Failed to add contact:', err);
	} finally {
		isAdding.value = false;
	}
}

// Remove a contact
async function removeContact(contact) {
	if (!confirm(`Remove ${contact.username} from your contacts?`)) {
		return;
	}

	try {
		await api.contacts.remove(props.userId, contact.id);

		// Remove from local list
		contacts.value = contacts.value.filter((c) => c.id !== contact.id);
	} catch (err) {
		error.value = 'Failed to remove contact';
		console.error('Failed to remove contact:', err);
	}
}

// Start conversation with contact
async function startConversation(contact) {
	try {
		const conversation = await api.conversations.create(props.userId, [
			props.userId,
			contact.contactUserId,
		]);

		// Emit event to parent to handle conversation creation
		emit('conversation-created', conversation);
	} catch (err) {
		error.value = 'Failed to start conversation';
		console.error('Failed to start conversation:', err);
	}
}

// Cancel add contact modal
function cancelAddContact() {
	showAddContact.value = false;
	searchQuery.value = '';
	searchResults.value = [];
	selectedUser.value = null;
	addError.value = '';
}

// Expose refresh method
defineExpose({
	refresh: loadContacts,
});
</script>

<style scoped>
.contacts-container {
	display: flex;
	flex-direction: column;
	height: 100%;
	padding: 16px;
}

.contacts-header {
	display: flex;
	justify-content: between;
	align-items: center;
	margin-bottom: 16px;
	padding-bottom: 12px;
	border-bottom: 1px solid var(--border-color);
}

.contacts-header h3 {
	margin: 0;
	flex: 1;
}

.contacts-loading,
.contacts-empty {
	text-align: center;
	padding: 32px 16px;
	color: var(--text-muted);
}

.contacts-list {
	flex: 1;
	overflow-y: auto;
}

.contact-item {
	display: flex;
	align-items: center;
	padding: 12px;
	border: 1px solid var(--border-color);
	border-radius: 8px;
	margin-bottom: 8px;
	background: var(--bg-surface);
}

.contact-item:hover {
	background: var(--hover-bg);
}

.contact-avatar {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	background: var(--avatar-bg);
	color: var(--avatar-text);
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: bold;
	margin-right: 12px;
}

.contact-info {
	flex: 1;
}

.contact-name {
	font-weight: 600;
	color: var(--text-primary);
}

.contact-id {
	font-size: 0.85em;
	color: var(--text-muted);
	font-family: monospace;
}

.contact-actions {
	display: flex;
	gap: 4px;
}

/* Modal styles */
.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
}

.add-contact-modal {
	background: var(--bg-surface);
	border-radius: 8px;
	width: 90%;
	max-width: 500px;
	max-height: 80vh;
	overflow-y: auto;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.modal-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 16px 20px;
	border-bottom: 1px solid var(--border-color);
}

.modal-header h3 {
	margin: 0;
}

.close-btn {
	background: none;
	border: none;
	font-size: 24px;
	cursor: pointer;
	color: var(--text-muted);
	padding: 0;
	width: 30px;
	height: 30px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.modal-body {
	padding: 20px;
}

.form-group {
	margin-bottom: 16px;
}

.form-group label {
	display: block;
	margin-bottom: 4px;
	font-weight: 600;
	color: var(--text-primary);
}

.form-control {
	width: 100%;
	padding: 8px 12px;
	border: 1px solid var(--border-color);
	border-radius: 4px;
	background: var(--bg-primary);
	color: var(--text-primary);
}

.search-results,
.selected-user {
	margin-top: 16px;
}

.search-results h4,
.selected-user h4 {
	margin: 0 0 8px 0;
	font-size: 16px;
	color: var(--text-primary);
}

.search-result-item,
.user-card {
	display: flex;
	align-items: center;
	padding: 8px;
	border: 1px solid var(--border-color);
	border-radius: 4px;
	margin-bottom: 4px;
	background: var(--bg-primary);
	cursor: pointer;
}

.search-result-item:hover {
	background: var(--hover-bg);
}

.user-card {
	cursor: default;
	background: var(--selected-bg);
	border-color: var(--selected-border);
}

.user-avatar {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	background: var(--avatar-bg);
	color: var(--avatar-text);
	display: flex;
	align-items: center;
	justify-content: center;
	font-weight: bold;
	margin-right: 8px;
	font-size: 14px;
}

.user-info {
	flex: 1;
}

.user-name {
	font-weight: 600;
	color: var(--text-primary);
}

.user-id {
	font-size: 0.8em;
	color: var(--text-muted);
	font-family: monospace;
}

.modal-footer {
	display: flex;
	justify-content: flex-end;
	gap: 8px;
	padding: 16px 20px;
	border-top: 1px solid var(--border-color);
}

.alert {
	padding: 8px 12px;
	border-radius: 4px;
	margin-bottom: 16px;
}

.alert-danger {
	background-color: var(--danger-bg);
	color: var(--danger-text);
	border: 1px solid var(--danger-text);
}

.btn {
	padding: 6px 12px;
	border: 1px solid transparent;
	border-radius: 4px;
	cursor: pointer;
	font-size: 14px;
	text-decoration: none;
	display: inline-block;
	text-align: center;
}

.btn-primary {
	background-color: var(--success-bg);
	color: var(--success-text);
	border-color: var(--success-bg);
}

.btn-secondary {
	background-color: var(--bg-secondary);
	color: var(--text-primary);
	border-color: var(--border-color);
}

.btn-outline-primary {
	background-color: transparent;
	color: var(--success-bg);
	border-color: var(--success-bg);
}

.btn-outline-danger {
	background-color: transparent;
	color: var(--danger-text);
	border-color: var(--danger-text);
}

.btn:hover {
	opacity: 0.9;
}

.btn:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}

.btn-sm {
	padding: 4px 8px;
	font-size: 12px;
}
</style>
