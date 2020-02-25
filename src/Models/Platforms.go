package Models

type PlatformConfig struct {
	Guid     string `json:"guid"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AddPlatformRequest struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EditPlatformRequest struct {
	Guid     string `json:"guid"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RemovePlatformRequest struct {
	Guid string `json:"guid"`
}
