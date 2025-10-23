<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			apiInfo: null,
		};
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				const api = await import('../services/api.js');
				this.apiInfo = await api.default.health.getInfo();
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		goToChat() {
			// Redirect to the main chat application
			this.$router.push('/');
		},
		exportList() {
			if (this.apiInfo) {
				const dataStr = JSON.stringify(this.apiInfo, null, 2);
				const dataBlob = new Blob([dataStr], {
					type: 'application/json',
				});
				const url = URL.createObjectURL(dataBlob);
				const link = document.createElement('a');
				link.href = url;
				link.download = 'wasa-api-info.json';
				link.click();
			}
		},
		newItem() {
			alert(
				'This is a demo page. The main chat application loads directly when you visit the root URL.'
			);
		},
	},
	mounted() {
		this.refresh();
	},
};
</script>

<template>
	<div class="container mt-4">
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
		>
			<h1 class="h2">WASAText API Information</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button
						type="button"
						class="btn btn-sm btn-outline-secondary"
						@click="refresh"
					>
						Refresh
					</button>
					<button
						type="button"
						class="btn btn-sm btn-outline-secondary"
						@click="exportList"
					>
						Export JSON
					</button>
				</div>
				<div class="btn-group me-2">
					<button
						type="button"
						class="btn btn-sm btn-outline-primary"
						@click="goToChat"
					>
						Go to Chat
					</button>
				</div>
			</div>
		</div>

		<LoadingSpinner v-if="loading" />
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div v-if="apiInfo && !loading" class="row">
			<div class="col-md-6">
				<div class="card">
					<div class="card-header">
						<h5 class="card-title">API Status</h5>
					</div>
					<div class="card-body">
						<p><strong>Name:</strong> {{ apiInfo.name }}</p>
						<p><strong>Version:</strong> {{ apiInfo.version }}</p>
						<p>
							<strong>Status:</strong>
							<span
								class="badge bg-success"
								v-if="apiInfo.status === 'running'"
								>{{ apiInfo.status }}</span
							>
							<span class="badge bg-warning" v-else>{{
								apiInfo.status
							}}</span>
						</p>
						<p>
							<strong>Available Endpoints:</strong>
							{{ apiInfo.endpoints || 'N/A' }}
						</p>
					</div>
				</div>
			</div>
			<div class="col-md-6">
				<div class="card">
					<div class="card-header">
						<h5 class="card-title">Features</h5>
					</div>
					<div class="card-body">
						<ul class="list-unstyled">
							<li>
								✅ Real-time messaging with WebSocket support
							</li>
							<li>
								✅ Image messaging with automatic compression
							</li>
							<li>✅ Emoji comments with toggle behavior</li>
							<li>✅ Group chat management</li>
							<li>✅ Contact management</li>
							<li>✅ Message forwarding (including images)</li>
							<li>✅ User profiles with photos</li>
						</ul>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style></style>
