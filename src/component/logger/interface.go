package logger

// Trace log category.
func Trace(category string, message string) {
	AddRuntimeLog(sessionID, "trace", category, message)
}

// Info log category.
func Info(category string, message string) {
	AddRuntimeLog(sessionID, "info", category, message)
}

// Debug log category.
func Debug(category string, message string) {
	AddRuntimeLog(sessionID, "debug", category, message)
}

// Warning log category.
func Warning(category string, message string) {
	AddRuntimeLog(sessionID, "warning", category, message)
}

// Error log category.
func Error(category string, message string) {
	AddRuntimeLog(sessionID, "error", category, message)
}

// Success log category.
func Success(category string, message string) {
	AddRuntimeLog(sessionID, "success", category, message)
}
