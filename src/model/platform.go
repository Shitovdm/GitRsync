package model

// PlatformConfig struct describes platform config.
type PlatformConfig struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// AddPlatformRequest struct describes add platform request model.
type AddPlatformRequest struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// EditPlatformRequest struct describes edit platform request model.
type EditPlatformRequest struct {
	UUID     string `json:"uuid"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// RemovePlatformRequest struct describes remove platform request model.
type RemovePlatformRequest struct {
	UUID string `json:"uuid"`
}
