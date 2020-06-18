package model

const (
	StateActive  = "active"
	StateBlocked = "blocked"
)

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

type AddRepositoryRequest struct {
	Name                    string `json:"name"`
	SourcePlatformUUID      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUUID string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
}

type EditRepositoryRequest struct {
	UUID                    string `json:"uuid"`
	Name                    string `json:"name"`
	SourcePlatformUUID      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUUID string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
}

type RemoveRepositoryRequest struct {
	UUID string `json:"uuid"`
}
