package Helpers

import (
	"errors"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
)

func GetAppConfig() (map[string]string, error) {
	var appConfig map[string]string
	conferr := Configuration.Load("AppConfig.json", &appConfig)
	if conferr != nil {
		return map[string]string{}, errors.New("unable to load app configuration")
	}

	return appConfig, nil
}

func GetRepositoriesConfig() ([]map[string]interface{}, error) {
	var repositoriesConfig []map[string]interface{}
	conferr := Configuration.Load("Repositories.json", &repositoriesConfig)
	if conferr != nil {
		return []map[string]interface{}{}, errors.New("unable to load repositories configuration")
	}

	return repositoriesConfig, nil
}

func GetPlatformsConfig() ([]map[string]interface{}, error) {
	var platformsConfig []map[string]interface{}
	conferr := Configuration.Load("Platforms.json", &platformsConfig)
	if conferr != nil {
		return []map[string]interface{}{}, errors.New("unable to load platforms configuration")
	}

	return platformsConfig, nil
}
