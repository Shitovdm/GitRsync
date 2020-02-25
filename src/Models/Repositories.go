package Models

type Repositories struct {
	Config []*RepositoryConfig `json:"config"`
}

type RepositoryConfig struct {
	Guid        string                   `json:"guid"`
	Source      RepositoryTransferConfig `json:"source"`
	Destination RepositoryTransferConfig `json:"destination"`
}

type RepositoryTransferConfig struct {
	PlatformGuid   string `json:"platformGuid"`
	RepositoryPath string `json:"repositoryPath"`
}
