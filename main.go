package main

import (
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Application"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Logger"
)

func init() {
	Configuration.Init("GitRsync")
}

func main() {

	//Helpers.Exec(`docker ps -a`)

	Logger.GetRuntimeLogFile()
	fmt.Println(Logger.GetRuntimeLogs())

	Application.StartServer()
}
