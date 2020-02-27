package main

import (
	"github.com/Shitovdm/git-repo-exporter/src/Application"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
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


	//Helpers.Exec(`docker ps -a`)

	/*repositories := Helpers.GetRepositoriesConfig()

	fmt.Println(repositories.Config)
	var platformsConfig map[string]interface{}
	_ = Configuration.Load("Platforms.json", &platformsConfig)
	fmt.Println(platformsConfig["config"])

	byteData, _ := json.Marshal(platformsConfig)
	fmt.Println(&byteData)*/
	Application.StartServer()

}
