package main

import (
	"github.com/Shitovdm/git-repo-exporter/src/Application"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

func init() {
	Configuration.Init("git-repo-exporter")
}

func main() {

	//systray.Run(onReady, onExit)

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


func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac.
	mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
}