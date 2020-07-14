package platform

import (
	"errors"
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/gofrs/uuid"
)

const (
	//	ConfigFileName describe platforms config file name.
	ConfigFileName = "Platforms.json"
)

// Platform struct describes platform config.
type Platform struct {
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

// GetUUID returns platform UUID.
func (p *Platform) GetUUID() string {
	return p.UUID
}

// SetUUID sets platform UUID.
func (p *Platform) SetUUID(UUID string) {
	p.UUID = UUID
}

// GetName returns platform Name.
func (p *Platform) GetName() string {
	return p.Name
}

// SetName sets platform Name.
func (p *Platform) SetName(name string) {
	p.Name = name
}

// GetAddress returns platform Address.
func (p *Platform) GetAddress() string {
	return p.Address
}

// SetAddress sets platform Address.
func (p *Platform) SetAddress(address string) {
	p.Address = address
}

// GetUsername returns platform Username.
func (p *Platform) GetUsername() string {
	return p.Username
}

// SetUsername sets platform Username.
func (p *Platform) SetUsername(username string) {
	p.Username = username
}

// GetPassword returns platform Password.
func (p *Platform) GetPassword() string {
	return p.Password
}

// SetPassword sets platform Password.
func (p *Platform) SetPassword(password string) {
	p.Password = password
}

// Get returns repository config by repository UUID.
func Get(UUID string) *Platform {

	platformsList := GetAll()
	for _, platform := range platformsList {
		if platform.UUID == UUID {
			return &platform
		}
	}

	return nil
}

// GetAll returns repositories config.
func GetAll() []Platform {

	platformsConfig := make([]Platform, 0)
	err := conf.Load(ConfigFileName, &platformsConfig)
	if err != nil {
		fmt.Printf("Error while loading platforms config file! %s", err.Error())
		err = conf.Save(ConfigFileName, []map[string]interface{}{})
		if err != nil {
			fmt.Printf("Error while creating new platforms config file! %s", err.Error())
		}
		return []Platform{}
	}

	return platformsConfig
}

// GetAllInInterface returns platform config data.
func GetAllInInterface() ([]map[string]interface{}, error) {

	var platformsConfig []map[string]interface{}
	err := conf.Load(ConfigFileName, &platformsConfig)
	if err != nil {
		return []map[string]interface{}{}, errors.New("unable to load platforms configuration")
	}
	return platformsConfig, nil
}

// Create creates new repository..
func Create(p *Platform) error {

	UUID, _ := uuid.NewV4()
	p.UUID = UUID.String()
	platforms := GetAll()
	platforms = append(platforms, *p)
	err := savePlatforms(platforms)
	if err != nil {
		return err
	}

	return nil
}

// Update describes update repository status action.
func (p *Platform) Update() error {

	oldPlatformsList := GetAll()
	for i, platform := range oldPlatformsList {
		if platform.UUID == p.UUID {
			oldPlatformsList[i].Name = p.Name
			oldPlatformsList[i].Address = p.Address
			oldPlatformsList[i].Username = p.Username
			oldPlatformsList[i].Password = p.Password
		}
	}

	err := savePlatforms(oldPlatformsList)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes repository config.
func (p *Platform) Delete() error {

	oldPlatformsList := GetAll()
	newPlatformsList := make([]Platform, 0)
	for _, platform := range oldPlatformsList {
		if platform.UUID != p.UUID {
			newPlatformsList = append(newPlatformsList, platform) //nolint:staticcheck
		}
	}

	err := savePlatforms(newPlatformsList)
	if err != nil {
		return err
	}

	return nil
}

// savePlatforms stores repositories config data.
func savePlatforms(repositories []Platform) error {
	err := conf.Save(ConfigFileName, &repositories)
	if err != nil {
		return fmt.Errorf("Error while saving repositories config file! %s ", err.Error())
	}
	return nil
}
