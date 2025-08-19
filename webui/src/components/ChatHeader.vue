<template>
	<div v-if="chat" class="chat-header d-flex align-items-center p-2 bg-white border-bottom shadow-sm">
		<div
			class="avatar bg-secondary rounded-circle me-2"
			style="width: 40px; height: 40px"
		></div>
		<div>
			<div class="fw-bold">{{ chat.name }}</div>
			<div class="text-muted small">
				{{ chat.online ? 'Online' : 'Offline' }}
			</div>
		</div>
	</div>
 </template>
