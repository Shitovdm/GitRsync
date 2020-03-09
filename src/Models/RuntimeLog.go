package Models

type RuntimeLog struct {
	SessionID string `json:"session_id"`
	Time      string `json:"time"`
	Level     string `json:"level"`
	Category  string `json:"category"`
	Message   string `json:"message"`
}

type ProcessLogsRequest struct {
	Action  string `json:"action"`
}