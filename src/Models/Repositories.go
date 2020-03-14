package Models

type RepositoryConfig struct {
	Uuid                    string `json:"uuid"`
	Name                    string `json:"name"`
	SourcePlatformUuid      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUuid string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
	Status                  string `json:"status"`
	UpdatedAt               string `json:"updated_at"`
}

type AddRepositoryRequest struct {
	Name                    string `json:"name"`
	SourcePlatformUuid      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUuid string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
}

type EditRepositoryRequest struct {
	Uuid                    string `json:"uuid"`
	Name                    string `json:"name"`
	SourcePlatformUuid      string `json:"spu"`
	SourcePlatformPath      string `json:"spp"`
	DestinationPlatformUuid string `json:"dpu"`
	DestinationPlatformPath string `json:"dpp"`
}

type RemoveRepositoryRequest struct {
	Uuid string `json:"uuid"`
}
