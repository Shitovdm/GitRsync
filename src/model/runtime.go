package model

// RuntimeLog struct describes runtime log config.
type RuntimeLog struct {
	SessionID string `json:"session_id"`
	Time      string `json:"time"`
	Level     string `json:"level"`
	Category  string `json:"category"`
	Message   string `json:"message"`
}

// RuntimeLogsRequest struct describes runtime log action.
type RuntimeLogsRequest struct {
	Action string `json:"action"`
}
