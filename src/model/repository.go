package model

const (
	// StateActive describe state "active"
	StateActive = "active"
	// StateBlocked describe state "blocked"
	StateBlocked = "blocked"
)

// RepositoryConfig struct describes repository config.
type RepositoryConfig struct {
	UUID                    string `json:"uuid"`
	Name                    string `json:"name"`
	SourcePlatformUUID      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUUID string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
	Status                  string `json:"status"`
	State                   string `json:"state"`
	UpdatedAt               string `json:"updated_at"`
}

// AddRepositoryRequest struct describes add repository request model.
type AddRepositoryRequest struct {
	Name                    string `json:"name"`
	SourcePlatformUUID      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUUID string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
}

// EditRepositoryRequest struct describes edit repository request model.
type EditRepositoryRequest struct {
	UUID                    string `json:"uuid"`
	Name                    string `json:"name"`
	SourcePlatformUUID      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUUID string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
}

// RemoveRepositoryRequest struct describes remove repository request model.
type RemoveRepositoryRequest struct {
	UUID string `json:"uuid"`
}
