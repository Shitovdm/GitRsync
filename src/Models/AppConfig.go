package Models

type AppConfig struct {
	Common            Common            `json:"common"`
	CommitsOverriding CommitsOverriding `json:"commits_overriding"`
}

type Common struct {
	RecentCommitsShown int `json:"recent_commits_shown"`
}

type CommitsOverriding struct {
	State           bool             `json:"state"`
	CommittersRules []CommittersRule `json:"committers_rules"`
}

type CommittersRule struct {
	Old GitUser `json:"old"`
	New GitUser `json:"new"`
}

type GitUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type SaveSettingsRequest struct {
	Section string      `json:"section"`
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
}
