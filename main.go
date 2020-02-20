package main

import (
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Application"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
)

func init() {
	Configuration.Init("git-repo-exporter")
}

func main() {

	/*appConfig := map[string]string{
		"gitlabLogin":    "test",
		"gitlabPassword": "test",
		"projectsFolder": "test",
		"shellstarter":   "test",
	}*/

	//_ = Configuration.Save("AppConfig.json", &appConfig)

	var newConfig interface{}
	res := Configuration.Load("AppConfig.json", &newConfig)
	if res == nil {
		fmt.Println(newConfig)
	}

	//Helpers.Exec(`docker ps -a`)

	fmt.Println(Helpers.GetCurrentPath())

	Application.StartServer()
}
