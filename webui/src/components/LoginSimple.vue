<template>
	<div class="login-container">
		<div class="login-card">
			<h1 class="login-title">WASAText - Debug Mode</h1>
			<p class="login-subtitle">Simple test version</p>

			<form @submit.prevent="handleLogin" class="login-form">
				<div class="input-group">
					<input
						v-model="username"
						type="text"
						placeholder="Username"
						class="login-input"
						required
					/>
				</div>

				<button type="submit" class="login-button">Test Login</button>
			</form>

			<div v-if="error" class="error-message">
				{{ error }}
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref } from 'vue';

console.log('Simple Login.vue script loaded');

const emit = defineEmits(['login-success']);

const username = ref('');
const error = ref('');

function handleLogin() {
	console.log('Simple handleLogin called with:', username.value);

	// Simulate successful login
	setTimeout(() => {
		console.log('Emitting fake login success');
		emit('login-success', {
			userId: 'test-user-id',
			username: username.value,
		});
	}, 1000);
}
</script>

<style scoped>
.login-container {
	display: flex;
	justify-content: center;
	align-items: center;
	min-height: 100vh;
	background: #f5f5f5;
	padding: 1rem;
}

.login-card {
	background: white;
	padding: 2rem;
	border-radius: 8px;
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
	width: 100%;
	max-width: 400px;
}

.login-title {
	text-align: center;
	margin-bottom: 1rem;
	color: #333;
}

.login-subtitle {
	text-align: center;
	margin-bottom: 2rem;
	color: #666;
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
	border: 1px solid #ddd;
	border-radius: 4px;
	font-size: 1rem;
}

.login-button {
	padding: 0.75rem;
	background: #007bff;
	color: white;
	border: none;
	border-radius: 4px;
	font-size: 1rem;
	cursor: pointer;
}

.login-button:hover {
	background: #0056b3;
}

.error-message {
	color: #dc3545;
	text-align: center;
	margin-top: 1rem;
}
</style>
