package Configuration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Shitovdm/git-repo-exporter/src/Models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type MarshalFunc func(v interface{}) ([]byte, error)

type UnmarshalFunc func(data []byte, v interface{}) error

var (
	applicationName = ""
	formats         = map[string]format{}
)

type format struct {
	m  MarshalFunc
	um UnmarshalFunc
}

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

func Init(application string) {
	applicationName = application
}

func Register(extension string, m MarshalFunc, um UnmarshalFunc) {
	formats[extension] = format{m, um}
}

func Load(path string, v interface{}) error {
	if applicationName == "" {
		panic("store: application name not defined")
	}

	if format, ok := formats[extension(path)]; ok {
		return LoadWith(path, v, format.um)
	}

	panic("store: unknown configuration format")
}

func Save(path string, v interface{}) error {
	if applicationName == "" {
		panic("store: application name not defined")
	}

	if format, ok := formats[extension(path)]; ok {
		return SaveWith(path, v, format.m)
	}

	panic("store: unknown configuration format")
}

func LoadWith(path string, v interface{}, um UnmarshalFunc) error {
	if applicationName == "" {
		panic("store: application name not defined")
	}

	globalPath := buildPlatformPath(path)

	data, err := ioutil.ReadFile(globalPath)

	if err != nil {
		return err
	}

	if err := um(data, v); err != nil {
		return fmt.Errorf("store: failed to unmarshal %s: %v", path, err)
	}

	return nil
}

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

	globalPath := buildPlatformPath(path)
	if err := os.MkdirAll(filepath.Dir(globalPath), os.ModePerm); err != nil {
		return err
	}

	if err := ioutil.WriteFile(globalPath, b.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func extension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}

	return ""
}

func buildPlatformPath(path string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s\\%s\\%s", os.Getenv("APPDATA"),
			applicationName,
			path)
	}

	var unixConfigDir string
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		unixConfigDir = xdg
	} else {
		unixConfigDir = os.Getenv("HOME") + "/.config"
	}

	return fmt.Sprintf("%s/%s/%s", unixConfigDir,
		applicationName,
		path)
}

func SetApplicationName(handle string) {
	applicationName = handle
}


func GetAppConfig() *Models.AppConfig {
	var appConfig Models.AppConfig
	err := Load("AppConfig.json", &appConfig)
	if err != nil {
		log.Println(fmt.Sprintf("Error while loading app config file! %s", err.Error()))
		err = Save("AppConfig.json", &Models.AppConfig{})
		if err != nil {
			log.Println(fmt.Sprintf("Error while creating new app config file! %s", err.Error()))
		}
		return &Models.AppConfig{}
	}

	return &appConfig
}

func GetRepositoriesConfig() []Models.RepositoryConfig {
	repositoriesConfig := make([]Models.RepositoryConfig, 0)
	err := Load("Repositories.json", &repositoriesConfig)
	if err != nil {
		log.Println(fmt.Sprintf("Error while loading repositories config file! %s", err.Error()))
		err = Save("Repositories.json", []map[string]interface{}{})
		if err != nil {
			log.Println(fmt.Sprintf("Error while creating new repositories config file! %s", err.Error()))
		}
		return []Models.RepositoryConfig{}
	}

	return repositoriesConfig
}

func GetRepositoriesConfigData() ([]map[string]interface{}, error) {

	var repositoriesConfig []map[string]interface{}
	err := Load("Repositories.json", &repositoriesConfig)
	if err != nil {
		return []map[string]interface{}{}, errors.New("unable to load repositories configuration")
	}
	return repositoriesConfig, nil
}

func SaveRepositoriesConfig(repositories []Models.RepositoryConfig) error {
	err := Save("Repositories.json", &repositories)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while saving repositories config file! %s", err.Error()))
	}
	return nil
}


func GetPlatformsConfig() []Models.PlatformConfig {
	platformsConfig := make([]Models.PlatformConfig, 0)
	err := Load("Platforms.json", &platformsConfig)
	if err != nil {
		log.Println(fmt.Sprintf("Error while loading platforms config file! %s", err.Error()))
		err = Save("Platforms.json", []map[string]interface{}{})
		if err != nil {
			log.Println(fmt.Sprintf("Error while creating new platforms config file! %s", err.Error()))
		}
		return []Models.PlatformConfig{}
	}

	return platformsConfig
}

func GetPlatformsConfigData() ([]map[string]interface{}, error) {

	var platformsConfig []map[string]interface{}
	err := Load("Platforms.json", &platformsConfig)
	if err != nil {
		return []map[string]interface{}{}, errors.New("unable to load platforms configuration")
	}
	return platformsConfig, nil
}

func SavePlatformsConfig(platforms []Models.PlatformConfig) error {
	err := Save("Platforms.json", &platforms)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while saving platforms config file! %s", err.Error()))
	}
	return nil
}