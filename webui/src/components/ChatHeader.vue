<template>
	<div
		v-if="chat"
		class="chat-header d-flex align-items-center p-2 border-bottom shadow-sm"
		style="
			background: var(--bg-surface);
			border-color: var(--border-color) !important;
		"
	>
		<div
			class="avatar rounded-circle me-2"
			style="width: 40px; height: 40px; background: var(--avatar-bg)"
		></div>
		<div>
			<div class="fw-bold" style="color: var(--text-primary)">
				{{ chat.name }}
			</div>
			<div class="small" style="color: var(--text-muted)">
				{{ chat.online ? 'Online' : 'Offline' }}
			</div>
		</div>
	</div>
</template>
