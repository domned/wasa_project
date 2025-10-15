package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// AdminHealthResponse represents the health status response
type AdminHealthResponse struct {
	Database  string `json:"database"`
	WebSocket string `json:"websocket"`
	API       string `json:"api"`
	Uptime    string `json:"uptime"`
}

// AdminStatsResponse represents the system statistics response
type AdminStatsResponse struct {
	TotalUsers         int     `json:"total_users"`
	ActiveUsers        int     `json:"active_users"`
	TotalConversations int     `json:"total_conversations"`
	TotalMessages      int     `json:"total_messages"`
	ActiveConnections  int     `json:"active_connections"`
	ErrorRate          float64 `json:"error_rate"`
}

// AdminLogsResponse represents the logs response
type AdminLogsResponse struct {
	Logs []database.LogEntry `json:"logs"`
}

var serverStartTime = time.Now()

// getAdminHealth returns the current system health status
func (rt *_router) getAdminHealth(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check database health
	databaseStatus := "Healthy"
	if err := rt.db.Ping(); err != nil {
		databaseStatus = "Error"
		rt.baseLogger.WithError(err).Error("Database ping failed")
	}

	// Check WebSocket health (count active connections)
	websocketStatus := "Active"
	activeConnections := 0
	if hub != nil {
		hub.mutex.RLock()
		activeConnections = len(hub.clients)
		hub.mutex.RUnlock()
		if activeConnections == 0 {
			websocketStatus = "Idle"
		}
	} else {
		websocketStatus = "Not initialized"
	}

	// Calculate uptime
	uptime := time.Since(serverStartTime)
	uptimeStr := formatUptime(uptime)

	response := AdminHealthResponse{
		Database:  databaseStatus,
		WebSocket: websocketStatus,
		API:       "Running",
		Uptime:    uptimeStr,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode admin health response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// getAdminStats returns system statistics
func (rt *_router) getAdminStats(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get total users count
	totalUsers, err := rt.db.GetUserCount()
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get user count")
		totalUsers = 0
	}

	// Get total conversations count
	totalConversations, err := rt.db.GetConversationCount()
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get conversation count")
		totalConversations = 0
	}

	// Get total messages count
	totalMessages, err := rt.db.GetMessageCount()
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get message count")
		totalMessages = 0
	}

	// Get active users count (last 24 hours)
	activeUsers, err := rt.db.GetActiveUserCount()
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get active user count")
		activeUsers = 0
	}

	// Get active WebSocket connections
	activeConnections := 0
	if hub != nil {
		hub.mutex.RLock()
		activeConnections = len(hub.clients)
		hub.mutex.RUnlock()
	}

	// Calculate error rate (placeholder - would need actual error tracking)
	errorRate := 0.01 // 1% default

	response := AdminStatsResponse{
		TotalUsers:         totalUsers,
		ActiveUsers:        activeUsers,
		TotalConversations: totalConversations,
		TotalMessages:      totalMessages,
		ActiveConnections:  activeConnections,
		ErrorRate:          errorRate,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode admin stats response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// getAdminLogs returns recent system logs
func (rt *_router) getAdminLogs(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get recent log entries from database
	logs, err := rt.db.GetRecentLogs(50) // Get last 50 log entries
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get recent logs")
		// Return empty logs on error
		logs = []database.LogEntry{}
	}

	response := AdminLogsResponse{
		Logs: logs,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode admin logs response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// formatUptime formats a duration into a human-readable uptime string
// getOnlineUsers returns list of currently online users based on WebSocket connections
func (rt *_router) getOnlineUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	onlineUserIDs := make(map[string]bool)

	// Get online users from WebSocket hub
	if hub != nil {
		hub.mutex.RLock()
		for client := range hub.clients {
			onlineUserIDs[client.UserID] = true
		}
		hub.mutex.RUnlock()
	}

	// Convert to slice
	var onlineUsers []string
	for userID := range onlineUserIDs {
		onlineUsers = append(onlineUsers, userID)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(onlineUsers); err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode online users response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func formatUptime(d time.Duration) string {
	if d < time.Minute {
		return "Less than a minute"
	}

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d hours, %d minutes", hours, minutes)
	} else {
		return fmt.Sprintf("%d minutes", minutes)
	}
}
