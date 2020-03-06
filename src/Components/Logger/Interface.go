package Logger

func Trace(category string, message string) {
	AddRuntimeLog(sessionID, "trace", category, message)
}

func Info(category string, message string) {
	AddRuntimeLog(sessionID, "info", category, message)
}

func Debug(category string, message string) {
	AddRuntimeLog(sessionID, "debug", category, message)
}

func Warning(category string, message string) {
	AddRuntimeLog(sessionID, "warning", category, message)
}

func Error(category string, message string) {
	AddRuntimeLog(sessionID, "error", category, message)
}
