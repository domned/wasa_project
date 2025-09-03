<template>
	<div class="login-container">
		<div class="login-card">
			<h1 class="login-title">WASAText</h1>
			<p class="login-subtitle">Enter your username to continue</p>

			<form @submit.prevent="handleLogin" class="login-form">
				<div class="input-group">
					<input
						v-model="username"
						type="text"
						placeholder="Username"
						class="login-input"
						:disabled="isLoading"
						required
						minlength="3"
						maxlength="16"
					/>
				</div>

				<button
					type="submit"
					class="login-button"
					:disabled="isLoading || username.length < 3"
				>
					<span v-if="isLoading">Logging in...</span>
					<span v-else>Login</span>
				</button>
			</form>

			<div v-if="error" class="error-message">
				{{ error }}
			</div>

			<div class="login-info">
				<p>
					<small
						>No password required. If the username doesn't exist, a
						new account will be created.</small
					>
				</p>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref } from 'vue';
import axios from '../services/axios.js';

const emit = defineEmits(['login-success']);

const username = ref('');
const isLoading = ref(false);
const error = ref('');

async function handleLogin() {
	if (username.value.length < 3 || username.value.length > 16) {
		error.value = 'Username must be between 3 and 16 characters';
		return;
	}

	isLoading.value = true;
	error.value = '';

	try {
		const response = await axios.post('/session', {
			name: username.value,
		});

		const userId = response.data.identifier;

		// Store user data in localStorage
		localStorage.setItem('userId', userId);
		localStorage.setItem('currentUsername', username.value);

		// Emit success event with user data
		emit('login-success', {
			userId: userId,
			username: username.value,
		});
	} catch (err) {
		console.error('Login error:', err);
		if (err.response && err.response.status === 400) {
			error.value = 'Invalid username. Please use 3-16 characters.';
		} else {
			error.value = 'Login failed. Please try again.';
		}
	} finally {
		isLoading.value = false;
	}
}
</script>

<style scoped>
.login-container {
	display: flex;
	justify-content: center;
	align-items: center;
	min-height: 100vh;
	background: var(--color-background);
	padding: 1rem;
}

.login-card {
	background: var(--color-surface);
	border-radius: 12px;
	padding: 2rem;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	width: 100%;
	max-width: 400px;
	border: 1px solid var(--color-border);
}

.login-title {
	font-size: 2rem;
	font-weight: 600;
	text-align: center;
	margin-bottom: 0.5rem;
	color: var(--color-text);
}

.login-subtitle {
	text-align: center;
	color: var(--color-text-secondary);
	margin-bottom: 2rem;
	font-size: 0.9rem;
}

.login-form {
	display: flex;
	flex-direction: column;
	gap: 1rem;
}

.input-group {
	display: flex;
	flex-direction: column;
}

.login-input {
	padding: 0.75rem;
	border: 1px solid var(--color-border);
	border-radius: 8px;
	font-size: 1rem;
	background: var(--color-background);
	color: var(--color-text);
	transition: border-color 0.2s ease;
}

.login-input:focus {
	outline: none;
	border-color: var(--color-primary);
	box-shadow: 0 0 0 3px rgba(13, 110, 253, 0.1);
}

.login-input:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}

.login-button {
	padding: 0.75rem;
	background: var(--color-primary);
	color: white;
	border: none;
	border-radius: 8px;
	font-size: 1rem;
	font-weight: 500;
	cursor: pointer;
	transition: all 0.2s ease;
}

.login-button:hover:not(:disabled) {
	background: var(--color-primary-hover);
	transform: translateY(-1px);
}

.login-button:active:not(:disabled) {
	transform: translateY(0);
}

.login-button:disabled {
	opacity: 0.6;
	cursor: not-allowed;
	transform: none;
}

.error-message {
	background: var(--color-danger-background);
	color: var(--color-danger);
	padding: 0.75rem;
	border-radius: 8px;
	text-align: center;
	margin-top: 1rem;
	border: 1px solid var(--color-danger-border);
}

.login-info {
	margin-top: 1.5rem;
	text-align: center;
}

.login-info small {
	color: var(--color-text-secondary);
	line-height: 1.4;
}

/* Dark theme support */
.dark .login-card {
	background: var(--color-surface);
	border-color: var(--color-border);
}

.dark .login-input {
	background: var(--color-background);
	border-color: var(--color-border);
	color: var(--color-text);
}

.dark .login-input:focus {
	border-color: var(--color-primary);
}
</style>
