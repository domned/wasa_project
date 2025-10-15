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

<script>
import { auth } from '../services/api.js';

export default {
	name: 'Login',
	emits: ['login-success'],
	data() {
		console.log('Login.vue data() called');
		return {
			username: '',
			isLoading: false,
			error: '',
		};
	},
	mounted() {
		console.log('Login component mounted - v2.0');
	},
	methods: {
		async handleLogin() {
			console.log('handleLogin function called!');
			if (this.username.length < 3 || this.username.length > 16) {
				this.error = 'Username must be between 3 and 16 characters';
				return;
			}

			this.isLoading = true;
			this.error = '';

			try {
				console.log('Attempting login with username:', this.username);

				// Make API call using the configured API service
				const data = await auth.login(this.username);
				const userId = data.identifier;
				console.log('Login successful, userId:', userId);

				// Store user data in localStorage
				localStorage.setItem('userId', userId);
				localStorage.setItem('currentUsername', this.username);
				console.log(
					'Stored in localStorage. UserId:',
					userId,
					'Username:',
					this.username
				);

				// Emit success event with user data
				console.log('Emitting login-success with:', {
					userId,
					username: this.username,
				});
				this.$emit('login-success', {
					userId: userId,
					username: this.username,
				});
			} catch (err) {
				console.error('Login error:', err);
				this.error = 'Login failed. Please try again.';
			} finally {
				this.isLoading = false;
			}
		},
	},
};
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
	background: #007bff !important;
	color: white !important;
	border: none;
	border-radius: 8px;
	font-size: 1rem;
	font-weight: 500;
	cursor: pointer;
	transition: all 0.2s ease;
	display: block !important;
	width: 100% !important;
	margin: 1rem 0 !important;
}

.login-button:hover:not(:disabled) {
	background: #0056b3 !important;
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
