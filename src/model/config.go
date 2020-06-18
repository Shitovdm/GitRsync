package model

// AppConfig struct describes app config.
type AppConfig struct {
	Common            Common            `json:"common"`
	CommitsOverriding CommitsOverriding `json:"commits_overriding"`
}

// Common struct describes app config common section.
type Common struct {
	RecentCommitsShown int `json:"recent_commits_shown"`
}

// CommitsOverriding struct describes app config commits overriding section.
type CommitsOverriding struct {
	State                        bool             `json:"state"`
	OverrideCommitsWithOneAuthor bool             `json:"override_commits_with_one_author"`
	MasterUser                   GitUser          `json:"master_user"`
	CommittersRules              []CommittersRule `json:"committers_rules"`
}

// CommittersRule struct describes committer rule info.
type CommittersRule struct {
	Old GitUser `json:"old"`
	New GitUser `json:"new"`
}

// GitUser struct describes git user info.
type GitUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// SaveSettingsRequest struct describes save app settings request.
type SaveSettingsRequest struct {
	Section string      `json:"section"`
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
}
