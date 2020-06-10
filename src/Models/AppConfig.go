package Models

type AppConfig struct {
	Common Common `json:"common"`
}

type Common struct {
	RecentCommitsShown int `json:"recent_commits_shown"`
}

type SaveSettingsRequest struct {
	Section string `json:"section"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}
