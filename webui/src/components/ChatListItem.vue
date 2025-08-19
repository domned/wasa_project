<template>
	<div
		class="chat-list-item d-flex align-items-center p-2 mb-1 rounded bg-white shadow-sm"
		style="cursor: pointer"
	>
		<div
			class="avatar bg-secondary rounded-circle me-2"
			style="width: 40px; height: 40px"
		></div>
		<div class="flex-grow-1">
			<div class="fw-bold">{{ chat.name }}</div>
			<div
				class="text-muted small text-truncate"
				style="max-width: 180px"
			>
				{{ chat.lastMessage }}
			</div>
		</div>
		<div class="text-end ms-2">
			<div class="small text-muted">{{ chat.time }}</div>
			<span v-if="chat.unread" class="badge bg-warning text-dark ms-1">{{
				chat.unread
			}}</span>
		</div>
	</div>
</template>

<script setup>
defineProps({
	chat: Object,
});
</script>

<style scoped>
.chat-list-item {
	transition: background 0.2s;
}
.chat-list-item:hover {
	background: #f0ede7;
}
.avatar {
	background: #b0b0b0;
}
</style>
