package Helpers

import (
	"errors"
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Models"
	"log"
)

func GetAppConfig() *Models.AppConfig {
	var appConfig Models.AppConfig
	err := Configuration.Load("AppConfig.json", &appConfig)
	if err != nil {
		log.Println(fmt.Sprintf("Error while loading app config file! %s", err.Error()))
		err = Configuration.Save("AppConfig.json", &Models.AppConfig{})
		if err != nil {
			log.Println(fmt.Sprintf("Error while creating new app config file! %s", err.Error()))
		}
		return &Models.AppConfig{}
	}

	return &appConfig
}

func GetRepositoriesConfig() *Models.Repositories {
	var repositoriesConfig Models.Repositories
	err := Configuration.Load("Repositories.json", &repositoriesConfig)
	if err != nil {
		log.Println(fmt.Sprintf("Error while loading repositories config file! %s", err.Error()))
		err = Configuration.Save("Repositories.json", &Models.Repositories{})
		if err != nil {
			log.Println(fmt.Sprintf("Error while creating new repositories config file! %s", err.Error()))
		}
		return &Models.Repositories{}
	}

	return &repositoriesConfig
}

func GetPlatformsConfig() []Models.PlatformConfig {
	platformsConfig := make([]Models.PlatformConfig, 0)

	err := Configuration.Load("Platforms.json", &platformsConfig)
	if err != nil {
		log.Println(fmt.Sprintf("Error while loading platforms config file! %s", err.Error()))
		err = Configuration.Save("Platforms.json", []map[string]interface{}{})
		if err != nil {
			log.Println(fmt.Sprintf("Error while creating new platforms config file! %s", err.Error()))
		}
		return []Models.PlatformConfig{}
	}

	return platformsConfig
}

func GetPlatformsConfigData() ([]map[string]interface{}, error) {

	var platformsConfig []map[string]interface{}
	err := Configuration.Load("Platforms.json", &platformsConfig)
	if err != nil {
		return []map[string]interface{}{}, errors.New("unable to load projects configuration")
	}
	return platformsConfig, nil
}

func SavePlatformsConfig(platforms []Models.PlatformConfig) error {
	err := Configuration.Save("Platforms.json", &platforms)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while saving platforms config file! %s", err.Error()))
	}
	return nil
}
