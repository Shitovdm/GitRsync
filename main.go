package main

import (
	"github.com/Shitovdm/git-repo-exporter/src/Application"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Logger"
)

func init() {
	Configuration.Init("GitRsync")
	Logger.Init()
}

func main() {

	//Helpers.Exec(`docker ps -a`)

	Application.StartServer()
}
