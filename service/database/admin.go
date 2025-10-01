package database

import (
	"time"
)

// GetUserCount returns the total number of users in the database
func (db *appdbimpl) GetUserCount() (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

// GetConversationCount returns the total number of conversations in the database
func (db *appdbimpl) GetConversationCount() (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM conversations").Scan(&count)
	return count, err
}

// GetMessageCount returns the total number of messages in the database
func (db *appdbimpl) GetMessageCount() (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	return count, err
}

// GetRecentLogs returns recent log entries from the database
func (db *appdbimpl) GetRecentLogs(limit int) ([]LogEntry, error) {
	query := `
		SELECT id, timestamp, level, message 
		FROM system_logs 
		ORDER BY timestamp DESC 
		LIMIT ?
	`
	
	rows, err := db.c.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		err := rows.Scan(&log.ID, &log.Timestamp, &log.Level, &log.Message)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// AddLogEntry adds a new log entry to the database
func (db *appdbimpl) AddLogEntry(level, message string) error {
	query := `
		INSERT INTO system_logs (timestamp, level, message) 
		VALUES (?, ?, ?)
	`
	
	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err := db.c.Exec(query, timestamp, level, message)
	return err
}

// UpdateUserLastSeen updates the last_seen timestamp for a user
func (db *appdbimpl) UpdateUserLastSeen(userID string) error {
	query := `UPDATE users SET last_seen = ? WHERE id = ?`
	timestamp := time.Now().Unix()
	_, err := db.c.Exec(query, timestamp, userID)
	return err
}

// GetActiveUserCount returns the number of users active in the last 24 hours
func (db *appdbimpl) GetActiveUserCount() (int, error) {
	var count int
	// Users active in the last 24 hours (86400 seconds)
	threshold := time.Now().Unix() - 86400
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE last_seen > ?", threshold).Scan(&count)
	return count, err
}