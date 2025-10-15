package api

import (
	"fmt"
	"time"
)

// SystemLogger provides structured logging to the database
type SystemLogger struct {
	rt *_router
}

// NewSystemLogger creates a new system logger instance
func NewSystemLogger(rt *_router) *SystemLogger {
	return &SystemLogger{rt: rt}
}

// LogInfo logs an info level message
func (sl *SystemLogger) LogInfo(message string) {
	sl.log("info", message)
}

// LogWarn logs a warning level message
func (sl *SystemLogger) LogWarn(message string) {
	sl.log("warn", message)
}

// LogError logs an error level message
func (sl *SystemLogger) LogError(message string) {
	sl.log("error", message)
}

// LogDebug logs a debug level message
func (sl *SystemLogger) LogDebug(message string) {
	sl.log("debug", message)
}

// log is the internal logging method
func (sl *SystemLogger) log(level, message string) {
	// Add timestamp to message for more context
	timestampedMessage := fmt.Sprintf("[%s] %s", time.Now().UTC().Format("15:04:05"), message)

	// Log to database
	if err := sl.rt.db.AddLogEntry(level, timestampedMessage); err != nil {
		sl.rt.baseLogger.WithError(err).Error("Failed to add log entry to database")
	}

	// Also log to application logger
	switch level {
	case "error":
		sl.rt.baseLogger.Error(message)
	case "warn":
		sl.rt.baseLogger.Warn(message)
	case "debug":
		sl.rt.baseLogger.Debug(message)
	default:
		sl.rt.baseLogger.Info(message)
	}
}
