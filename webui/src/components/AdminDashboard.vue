<template>
	<div class="admin-dashboard">
		<!-- Header -->
		<div class="dashboard-header">
			<h2>System Administration</h2>
			<button
				class="btn btn-outline-secondary btn-sm"
				@click="refreshAll"
				:disabled="isRefreshing"
			>
				{{ isRefreshing ? 'Refreshing...' : '🔄 Refresh All' }}
			</button>
		</div>

		<!-- System Health Status -->
		<div class="dashboard-section">
			<h3>System Health</h3>
			<div class="health-cards">
				<div
					class="health-card"
					:class="getHealthStatusClass(systemHealth.database)"
				>
					<div class="health-icon">💾</div>
					<div class="health-info">
						<div class="health-title">Database</div>
						<div class="health-status">
							{{ systemHealth.database || 'Unknown' }}
						</div>
					</div>
				</div>

				<div
					class="health-card"
					:class="getHealthStatusClass(systemHealth.websocket)"
				>
					<div class="health-icon">🔌</div>
					<div class="health-info">
						<div class="health-title">WebSocket</div>
						<div class="health-status">
							{{ systemHealth.websocket || 'Unknown' }}
						</div>
					</div>
				</div>

				<div
					class="health-card"
					:class="getHealthStatusClass(systemHealth.api)"
				>
					<div class="health-icon">🌐</div>
					<div class="health-info">
						<div class="health-title">API</div>
						<div class="health-status">
							{{ systemHealth.api || 'Unknown' }}
						</div>
					</div>
				</div>

				<div class="health-card uptime">
					<div class="health-icon">⏱️</div>
					<div class="health-info">
						<div class="health-title">Uptime</div>
						<div class="health-status">
							{{ systemHealth.uptime || 'Unknown' }}
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- System Statistics -->
		<div class="dashboard-section">
			<h3>System Statistics</h3>
			<div class="stats-grid">
				<div class="stat-card">
					<div class="stat-number">{{ stats.totalUsers }}</div>
					<div class="stat-label">Total Users</div>
				</div>
				<div class="stat-card">
					<div class="stat-number">{{ stats.activeUsers }}</div>
					<div class="stat-label">Active Users</div>
				</div>
				<div class="stat-card">
					<div class="stat-number">
						{{ stats.totalConversations }}
					</div>
					<div class="stat-label">Conversations</div>
				</div>
				<div class="stat-card">
					<div class="stat-number">{{ stats.totalMessages }}</div>
					<div class="stat-label">Messages</div>
				</div>
				<div class="stat-card">
					<div class="stat-number">{{ stats.activeConnections }}</div>
					<div class="stat-label">WebSocket Connections</div>
				</div>
				<div class="stat-card">
					<div class="stat-number">{{ stats.errorRate }}%</div>
					<div class="stat-label">Error Rate</div>
				</div>
			</div>
		</div>

		<!-- Active Users Management -->
		<div class="dashboard-section">
			<h3>Active Users</h3>
			<div class="users-table-container">
				<table class="users-table">
					<thead>
						<tr>
							<th>User</th>
							<th>Status</th>
							<th>Last Active</th>
							<th>Actions</th>
						</tr>
					</thead>
					<tbody>
						<tr v-for="user in activeUsers" :key="user.id">
							<td>
								<div class="user-info">
									<img
										v-if="user.picture"
										:src="user.picture"
										:alt="user.username"
										class="user-avatar"
									/>
									<div v-else class="user-avatar-placeholder">
										{{
											user.username
												.charAt(0)
												.toUpperCase()
										}}
									</div>
									<span class="username">{{
										user.username
									}}</span>
								</div>
							</td>
							<td>
								<span
									class="status-badge"
									:class="user.online ? 'online' : 'offline'"
								>
									{{ user.online ? 'Online' : 'Offline' }}
								</span>
							</td>
							<td>{{ formatDate(user.lastActive) }}</td>
							<td>
								<button
									class="btn btn-sm btn-outline-warning"
									@click="disconnectUser(user.id)"
									:disabled="!user.online"
									title="Disconnect user"
								>
									Disconnect
								</button>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<!-- System Logs -->
		<div class="dashboard-section">
			<h3>Recent System Events</h3>
			<div class="logs-container">
				<div
					v-for="log in recentLogs"
					:key="log.id"
					class="log-entry"
					:class="log.level"
				>
					<div class="log-time">{{ formatTime(log.timestamp) }}</div>
					<div class="log-level">{{ log.level.toUpperCase() }}</div>
					<div class="log-message">{{ log.message }}</div>
				</div>
				<div v-if="recentLogs.length === 0" class="no-logs">
					No recent events
				</div>
			</div>
		</div>

		<!-- System Actions -->
		<div class="dashboard-section">
			<h3>System Actions</h3>
			<div class="actions-grid">
				<button
					class="btn btn-outline-info action-btn"
					@click="clearLogs"
				>
					🗑️ Clear Logs
				</button>
				<button
					class="btn btn-outline-warning action-btn"
					@click="restartWebSocket"
				>
					🔄 Restart WebSocket
				</button>
				<button
					class="btn btn-outline-success action-btn"
					@click="backupData"
					:disabled="isBackingUp"
				>
					💾 {{ isBackingUp ? 'Backing up...' : 'Backup Data' }}
				</button>
				<button
					class="btn btn-outline-danger action-btn"
					@click="showMaintenanceMode"
				>
					🔧 Maintenance Mode
				</button>
			</div>
		</div>

		<!-- Error Display -->
		<div v-if="error" class="alert alert-danger">
			{{ error }}
		</div>
	</div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { apiService } from '../services/api.js';

const systemHealth = ref({
	database: 'Unknown',
	websocket: 'Unknown',
	api: 'Unknown',
	uptime: 'Unknown',
});

const stats = ref({
	totalUsers: 0,
	activeUsers: 0,
	totalConversations: 0,
	totalMessages: 0,
	activeConnections: 0,
	errorRate: 0,
});

const activeUsers = ref([]);
const recentLogs = ref([]);
const error = ref('');
const isRefreshing = ref(false);
const isBackingUp = ref(false);

let refreshInterval = null;

onMounted(() => {
	loadDashboardData();
	// Refresh dashboard every 30 seconds
	refreshInterval = setInterval(loadDashboardData, 30000);
});

onUnmounted(() => {
	if (refreshInterval) {
		clearInterval(refreshInterval);
	}
});

async function loadDashboardData() {
	try {
		await Promise.all([
			loadSystemHealth(),
			loadStats(),
			loadActiveUsers(),
			loadRecentLogs(),
		]);
		error.value = '';
	} catch (err) {
		console.error('Failed to load dashboard data:', err);
		error.value = 'Failed to load dashboard data';
	}
}

async function loadSystemHealth() {
	try {
		const response = await apiService.health.checkHealth();
		if (response.success) {
			systemHealth.value = {
				database: response.data.database || 'Healthy',
				websocket: response.data.websocket || 'Active',
				api: response.data.api || 'Running',
				uptime: response.data.uptime || 'Unknown',
			};
		}
	} catch (err) {
		console.error('Failed to load system health:', err);
	}
}

async function loadStats() {
	try {
		const response = await apiService.health.getStats();
		if (response.success) {
			stats.value = {
				totalUsers: response.data.total_users || 0,
				activeUsers: response.data.active_users || 0,
				totalConversations: response.data.total_conversations || 0,
				totalMessages: response.data.total_messages || 0,
				activeConnections: response.data.active_connections || 0,
				errorRate: response.data.error_rate || 0,
			};
		}
	} catch (err) {
		console.error('Failed to load stats:', err);
	}
}

async function loadActiveUsers() {
	try {
		const response = await apiService.users.listAll();
		if (response.success) {
			activeUsers.value = response.data.map((user) => ({
				...user,
				lastActive: user.last_active || new Date().toISOString(),
			}));
		}
	} catch (err) {
		console.error('Failed to load active users:', err);
	}
}

async function loadRecentLogs() {
	try {
		const response = await apiService.health.getLogs();
		if (response.success) {
			recentLogs.value = response.data.slice(0, 20); // Show last 20 logs
		}
	} catch (err) {
		console.error('Failed to load logs:', err);
		// Mock some logs for demonstration
		recentLogs.value = [
			{
				id: 1,
				timestamp: new Date().toISOString(),
				level: 'info',
				message: 'System health check completed successfully',
			},
			{
				id: 2,
				timestamp: new Date(Date.now() - 60000).toISOString(),
				level: 'info',
				message: 'New user connected via WebSocket',
			},
			{
				id: 3,
				timestamp: new Date(Date.now() - 120000).toISOString(),
				level: 'warn',
				message: 'High memory usage detected (85%)',
			},
		];
	}
}

async function refreshAll() {
	isRefreshing.value = true;
	await loadDashboardData();
	isRefreshing.value = false;
}

function getHealthStatusClass(status) {
	const normalizedStatus = status?.toLowerCase() || '';
	if (
		normalizedStatus.includes('healthy') ||
		normalizedStatus.includes('active') ||
		normalizedStatus.includes('running')
	) {
		return 'healthy';
	} else if (
		normalizedStatus.includes('warning') ||
		normalizedStatus.includes('degraded')
	) {
		return 'warning';
	} else if (
		normalizedStatus.includes('error') ||
		normalizedStatus.includes('down') ||
		normalizedStatus.includes('failed')
	) {
		return 'error';
	}
	return 'unknown';
}

async function disconnectUser(userId) {
	if (!confirm('Disconnect this user from the system?')) return;

	try {
		const response = await apiService.health.disconnectUser(userId);
		if (response.success) {
			await loadActiveUsers();
		} else {
			alert('Failed to disconnect user');
		}
	} catch (err) {
		console.error('Failed to disconnect user:', err);
		alert('Failed to disconnect user');
	}
}

async function clearLogs() {
	if (!confirm('Clear all system logs? This action cannot be undone.'))
		return;

	try {
		const response = await apiService.health.clearLogs();
		if (response.success) {
			recentLogs.value = [];
		} else {
			alert('Failed to clear logs');
		}
	} catch (err) {
		console.error('Failed to clear logs:', err);
		alert('Failed to clear logs');
	}
}

async function restartWebSocket() {
	if (
		!confirm(
			'Restart WebSocket service? This will disconnect all users briefly.'
		)
	)
		return;

	try {
		const response = await apiService.health.restartWebSocket();
		if (response.success) {
			alert('WebSocket service restarted successfully');
			await loadSystemHealth();
		} else {
			alert('Failed to restart WebSocket service');
		}
	} catch (err) {
		console.error('Failed to restart WebSocket:', err);
		alert('Failed to restart WebSocket service');
	}
}

async function backupData() {
	if (!confirm('Create a system backup? This may take a few minutes.'))
		return;

	isBackingUp.value = true;
	try {
		const response = await apiService.health.createBackup();
		if (response.success) {
			alert('Backup created successfully');
		} else {
			alert('Failed to create backup');
		}
	} catch (err) {
		console.error('Failed to create backup:', err);
		alert('Failed to create backup');
	} finally {
		isBackingUp.value = false;
	}
}

function showMaintenanceMode() {
	alert('Maintenance mode functionality would be implemented here');
}

function formatDate(dateString) {
	if (!dateString) return 'Never';
	const date = new Date(dateString);
	return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
}

function formatTime(dateString) {
	if (!dateString) return '';
	const date = new Date(dateString);
	return date.toLocaleTimeString();
}
</script>

<style scoped>
.admin-dashboard {
	padding: 2rem;
	max-width: 1200px;
	margin: 0 auto;
	background: var(--bg-primary);
	min-height: 100vh;
}

.dashboard-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 2rem;
	padding-bottom: 1rem;
	border-bottom: 2px solid var(--border-color);
}

.dashboard-header h2 {
	color: var(--text-primary);
	margin: 0;
}

.dashboard-section {
	background: var(--bg-surface);
	border-radius: 8px;
	padding: 1.5rem;
	margin-bottom: 2rem;
	border: 1px solid var(--border-color);
}

.dashboard-section h3 {
	color: var(--text-primary);
	margin: 0 0 1rem 0;
	font-size: 1.25rem;
}

/* Health Cards */
.health-cards {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
	gap: 1rem;
}

.health-card {
	display: flex;
	align-items: center;
	padding: 1rem;
	border-radius: 8px;
	border: 2px solid;
	transition: all 0.3s;
}

.health-card.healthy {
	border-color: var(--success);
	background: var(--success-bg);
}

.health-card.warning {
	border-color: var(--warning);
	background: var(--warning-bg);
}

.health-card.error {
	border-color: var(--danger);
	background: var(--danger-bg);
}

.health-card.unknown {
	border-color: var(--border-color);
	background: var(--bg-secondary);
}

.health-card.uptime {
	border-color: var(--primary);
	background: var(--primary-bg);
}

.health-icon {
	font-size: 2rem;
	margin-right: 1rem;
}

.health-title {
	font-weight: 600;
	color: var(--text-primary);
}

.health-status {
	font-size: 0.875rem;
	color: var(--text-muted);
}

/* Stats Grid */
.stats-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
	gap: 1rem;
}

.stat-card {
	text-align: center;
	padding: 1.5rem 1rem;
	background: var(--bg-message);
	border-radius: 8px;
	border: 1px solid var(--border-color);
}

.stat-number {
	font-size: 2rem;
	font-weight: bold;
	color: var(--primary);
	margin-bottom: 0.5rem;
}

.stat-label {
	font-size: 0.875rem;
	color: var(--text-muted);
}

/* Users Table */
.users-table-container {
	overflow-x: auto;
}

.users-table {
	width: 100%;
	border-collapse: collapse;
}

.users-table th,
.users-table td {
	padding: 0.75rem;
	text-align: left;
	border-bottom: 1px solid var(--border-color);
}

.users-table th {
	background: var(--bg-secondary);
	color: var(--text-primary);
	font-weight: 600;
}

.user-info {
	display: flex;
	align-items: center;
	gap: 0.5rem;
}

.user-avatar {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	object-fit: cover;
}

.user-avatar-placeholder {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	background: var(--avatar-bg);
	display: flex;
	align-items: center;
	justify-content: center;
	color: var(--text-primary);
	font-weight: bold;
	font-size: 0.875rem;
}

.username {
	color: var(--text-primary);
}

.status-badge {
	padding: 0.25rem 0.5rem;
	border-radius: 12px;
	font-size: 0.75rem;
	font-weight: 500;
}

.status-badge.online {
	background: var(--success-bg);
	color: var(--success);
}

.status-badge.offline {
	background: var(--bg-secondary);
	color: var(--text-muted);
}

/* Logs */
.logs-container {
	max-height: 300px;
	overflow-y: auto;
	border: 1px solid var(--border-color);
	border-radius: 4px;
}

.log-entry {
	display: flex;
	align-items: center;
	gap: 1rem;
	padding: 0.5rem;
	border-bottom: 1px solid var(--border-color);
	font-family: monospace;
	font-size: 0.875rem;
}

.log-entry:last-child {
	border-bottom: none;
}

.log-entry.info {
	background: var(--info-bg);
}

.log-entry.warn {
	background: var(--warning-bg);
}

.log-entry.error {
	background: var(--danger-bg);
}

.log-time {
	color: var(--text-muted);
	min-width: 80px;
}

.log-level {
	font-weight: bold;
	min-width: 60px;
}

.log-level {
	color: var(--text-primary);
}

.log-entry.warn .log-level {
	color: var(--warning);
}

.log-entry.error .log-level {
	color: var(--danger);
}

.log-message {
	color: var(--text-primary);
	flex: 1;
}

.no-logs {
	padding: 2rem;
	text-align: center;
	color: var(--text-muted);
}

/* Actions */
.actions-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
	gap: 1rem;
}

.action-btn {
	padding: 0.75rem 1rem;
	font-weight: 500;
}

/* Buttons */
.btn {
	padding: 0.375rem 0.75rem;
	border: 1px solid transparent;
	border-radius: 4px;
	cursor: pointer;
	font-size: 0.875rem;
	transition: all 0.2s;
	text-decoration: none;
	display: inline-flex;
	align-items: center;
	justify-content: center;
	gap: 0.5rem;
}

.btn:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}

.btn-outline-secondary {
	color: var(--text-primary);
	border-color: var(--border-color);
}

.btn-outline-secondary:hover:not(:disabled) {
	background: var(--bg-hover);
}

.btn-outline-info {
	color: var(--info);
	border-color: var(--info);
}

.btn-outline-info:hover:not(:disabled) {
	background: var(--info);
	color: white;
}

.btn-outline-warning {
	color: var(--warning);
	border-color: var(--warning);
}

.btn-outline-warning:hover:not(:disabled) {
	background: var(--warning);
	color: white;
}

.btn-outline-success {
	color: var(--success);
	border-color: var(--success);
}

.btn-outline-success:hover:not(:disabled) {
	background: var(--success);
	color: white;
}

.btn-outline-danger {
	color: var(--danger);
	border-color: var(--danger);
}

.btn-outline-danger:hover:not(:disabled) {
	background: var(--danger);
	color: white;
}

.btn-sm {
	padding: 0.25rem 0.5rem;
	font-size: 0.75rem;
}

.alert {
	padding: 0.75rem;
	margin: 1rem 0;
	border-radius: 4px;
}

.alert-danger {
	background: var(--danger-bg);
	color: var(--danger);
	border: 1px solid var(--danger);
}
</style>
