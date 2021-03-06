package conf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Shitovdm/GitRsync/src/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// MarshalFunc returns bytes array.
type MarshalFunc func(v interface{}) ([]byte, error)

// UnmarshalFunc gets bytes array.
type UnmarshalFunc func(data []byte, v interface{}) error

var (
	applicationName = ""
	formats         = map[string]format{}
)

// format struct describes config format.
type format struct {
	m  MarshalFunc
	um UnmarshalFunc
}

// init describes config format.
func init() {
	formats["json"] = format{m: json.Marshal, um: json.Unmarshal}
	formats["yaml"] = format{m: yaml.Marshal, um: yaml.Unmarshal}
	formats["yml"] = format{m: yaml.Marshal, um: yaml.Unmarshal}

	formats["toml"] = format{
		m: func(v interface{}) ([]byte, error) {
			b := bytes.Buffer{}
			err := toml.NewEncoder(&b).Encode(v)
			return b.Bytes(), err
		},
		um: toml.Unmarshal,
	}
}

// Init sets application name.
func Init(application string) {
	applicationName = application
}

// Load gets config file.
func Load(path string, v interface{}) error {
	if applicationName == "" {
		panic("store: application name not defined")
	}

	if format, ok := formats[extension(path)]; ok {
		return LoadWith(path, v, format.um)
	}

	panic("store: unknown configuration format")
}

// Save stores file.
func Save(path string, v interface{}) error {
	if applicationName == "" {
		panic("store: application name not defined")
	}

	if format, ok := formats[extension(path)]; ok {
		return SaveWith(path, v, format.m)
	}

	panic("store: unknown configuration format")
}

// LoadWith gets file.
func LoadWith(path string, v interface{}, um UnmarshalFunc) error {
	if applicationName == "" {
		panic("store: application name not defined")
	}

	globalPath := BuildPlatformPath(path)

	data, err := ioutil.ReadFile(globalPath)

	if err != nil {
		return err
	}

	if err := um(data, v); err != nil {
		return fmt.Errorf("store: failed to unmarshal %s: %v", path, err)
	}

	return nil
}

// SaveWith stores file.
func SaveWith(path string, v interface{}, m MarshalFunc) error {
	if applicationName == "" {
		panic("store: application name not defined")
	}

	var b bytes.Buffer

	if data, err := m(v); err == nil {
		b.Write(data)
	} else {
		return fmt.Errorf("store: failed to marshal %s: %v", path, err)
	}

	b.WriteRune('\n')

	globalPath := BuildPlatformPath(path)
	if err := os.MkdirAll(filepath.Dir(globalPath), os.ModePerm); err != nil {
		return err
	}

	if err := ioutil.WriteFile(globalPath, b.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return nil
}

// extension returns filepath extension.
func extension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}

	return ""
}

// BuildPlatformPath returns platform path.
func BuildPlatformPath(path string) string {

	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s\\%s\\%s", os.Getenv("APPDATA"), GetApplicationName(), path)
	}

	var unixConfigDir string
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		unixConfigDir = xdg
	} else {
		unixConfigDir = os.Getenv("HOME") + "/.config"
	}

	return fmt.Sprintf("%s/%s/%s", unixConfigDir, GetApplicationName(), path)
}

// GetApplicationName returns application name.
func GetApplicationName() string {
	return applicationName
}

// GetAppConfig gets app config model.
func GetAppConfig() *model.AppConfig {
	var appConfig model.AppConfig
	err := Load("AppConfig.json", &appConfig)
	if err != nil {
		fmt.Printf("Error while loading app config file! %s", err.Error())
		err = Save("AppConfig.json", &model.AppConfig{})
		if err != nil {
			fmt.Printf("Error while creating new app config file! %s", err.Error())
		}
		return &model.AppConfig{}
	}

	return &appConfig
}

// GetAppConfigData returns app config data.
func GetAppConfigData() (map[string]interface{}, error) {

	var appConfig map[string]interface{}
	err := Load("AppConfig.json", &appConfig)
	if err != nil {
		return map[string]interface{}{}, errors.New("unable to load app configuration")
	}
	return appConfig, nil
}

// GetAppConfigField returns app config field.
func GetAppConfigField(section string, field string) reflect.Value {

	var appConfig model.AppConfig
	_ = Load("AppConfig.json", &appConfig)
	return reflect.Indirect(reflect.ValueOf(appConfig)).FieldByName(section).FieldByName(field)
}

// SaveAppConfig stores app config.
func SaveAppConfig(appConfig *model.AppConfig) error {

	err := Save("AppConfig.json", &appConfig)
	if err != nil {
		return fmt.Errorf("Error while saving app config file! %s ", err.Error())
	}
	return nil
}

// GetRepositoriesConfig returns repositories config.
func GetRepositoriesConfig() []model.RepositoryConfig {

	repositoriesConfig := make([]model.RepositoryConfig, 0)
	err := Load("Repositories.json", &repositoriesConfig)
	if err != nil {
		fmt.Printf("Error while loading repositories config file! %s", err.Error())
		err = Save("Repositories.json", []map[string]interface{}{})
		if err != nil {
			fmt.Printf("Error while creating new repositories config file! %s", err.Error())
		}
		return []model.RepositoryConfig{}
	}

	return repositoriesConfig
}

// GetRepositoriesConfigData returns repositories config data.
func GetRepositoriesConfigData() ([]map[string]interface{}, error) {

	var repositoriesConfig []map[string]interface{}
	err := Load("Repositories.json", &repositoriesConfig)
	if err != nil {
		return []map[string]interface{}{}, errors.New("unable to load repositories configuration")
	}
	return repositoriesConfig, nil
}

// GetActiveRepositoriesConfigData returns active repositories config.
func GetActiveRepositoriesConfigData() ([]map[string]interface{}, error) {

	repositoriesConfig, _ := GetRepositoriesConfigData()
	var activeRepositories []map[string]interface{}
	for _, repo := range repositoriesConfig {
		if repo["state"] == model.StateActive {
			activeRepositories = append(activeRepositories, repo)
		}
	}

	return activeRepositories, nil
}

// GetBlockedRepositoriesConfigData returns blocked repositories config.
func GetBlockedRepositoriesConfigData() ([]map[string]interface{}, error) {

	repositoriesConfig, _ := GetRepositoriesConfigData()
	var blockedRepositories []map[string]interface{}
	for _, repo := range repositoriesConfig {
		if repo["state"] == model.StateBlocked {
			blockedRepositories = append(blockedRepositories, repo)
		}
	}

	return blockedRepositories, nil
}

// GetRepositoryDestinationRepositoryName parses destination repository name.
func GetRepositoryDestinationRepositoryName(repoConfig *model.RepositoryConfig) string {

	dpp := strings.Trim(strings.TrimRight(repoConfig.DestinationPlatformPath, "git"), ".")
	return strings.Split(dpp, "/")[len(strings.Split(dpp, "/"))-1]
}

// GetRepositorySourceRepositoryName parses source repository name.
func GetRepositorySourceRepositoryName(repoConfig *model.RepositoryConfig) string {

	spp := strings.Trim(strings.TrimRight(repoConfig.SourcePlatformPath, "git"), ".")
	return strings.Split(spp, "/")[len(strings.Split(spp, "/"))-1]
}

// SaveRepositoriesConfig stores repositories config.
func SaveRepositoriesConfig(repositories []model.RepositoryConfig) error {

	err := Save("Repositories.json", &repositories)
	if err != nil {
		return fmt.Errorf("Error while saving repositories config file! %s ", err.Error())
	}
	return nil
}

// GetPlatformsConfig gets platforms config.
func GetPlatformsConfig() []model.PlatformConfig {

	platformsConfig := make([]model.PlatformConfig, 0)
	err := Load("Platforms.json", &platformsConfig)
	if err != nil {
		fmt.Printf("Error while loading platforms config file! %s", err.Error())
		err = Save("Platforms.json", []map[string]interface{}{})
		if err != nil {
			fmt.Printf("Error while creating new platforms config file! %s", err.Error())
		}
		return []model.PlatformConfig{}
	}

	return platformsConfig
}

// GetPlatformsConfigData returns platform config data.
func GetPlatformsConfigData() ([]map[string]interface{}, error) {

	var platformsConfig []map[string]interface{}
	err := Load("Platforms.json", &platformsConfig)
	if err != nil {
		return []map[string]interface{}{}, errors.New("unable to load platforms configuration")
	}
	return platformsConfig, nil
}

// SavePlatformsConfig stores platform config.
func SavePlatformsConfig(platforms []model.PlatformConfig) error {

	err := Save("Platforms.json", &platforms)
	if err != nil {
		return fmt.Errorf("Error while saving platforms config file! %s ", err.Error())
	}
	return nil
}

// GetPlatformByUUID returns platform config by platform UUID.
func GetPlatformByUUID(UUID string) *model.PlatformConfig {

	platformsList := GetPlatformsConfig()
	for _, platform := range platformsList {
		if platform.UUID == UUID {
			return &platform
		}
	}

	return nil
}

// GetRepositoryByUUID returns repository config by repository UUID.
func GetRepositoryByUUID(UUID string) *model.RepositoryConfig {

	repositoriesList := GetRepositoriesConfig()
	for _, repository := range repositoriesList {
		if repository.UUID == UUID {
			return &repository
		}
	}

	return nil
}
