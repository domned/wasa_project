package api

// SystemLogger provides structured logging
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
	// Log to application logger only (database logging removed)
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
