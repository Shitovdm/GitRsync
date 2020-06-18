package model

type PlatformConfig struct {
	UUID     string `json:"uuid"`
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
	UUID     string `json:"uuid"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RemovePlatformRequest struct {
	UUID string `json:"uuid"`
}
