<template>
	<div class="login-container">
		<div class="login-card">
			<h1 class="login-title">WASAText - DEBUG</h1>
			<p class="login-subtitle">Simple test version</p>
			<button @click="testLogin" class="login-button">TEST LOGIN</button>
			<div>Debug: Component loaded successfully</div>
		</div>
	</div>
</template>

<script setup>
console.log('Simple Login component loaded');

function testLogin() {
	console.log('Test login clicked');
	alert('Test login works!');
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
	text-align: center;
}

.login-button {
	background: #007bff;
	color: white;
	border: none;
	padding: 0.75rem 1.5rem;
	border-radius: 4px;
	cursor: pointer;
	margin: 1rem 0;
}
</style>
