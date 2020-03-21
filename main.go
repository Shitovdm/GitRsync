package main

import (
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Cmd"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Logger"
	"time"
)

func init() {
	Configuration.Init("GitRsync")
	Logger.Init()
}

func main() {

	resource := "https://github.com/Shitovdm/rpc.git"
	path := Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s", "rpc-test-sync"))


	/*commits, _ := Cmd.Log(path, "ba3edcc592c66d40b18613ac044d0bcf277eb08b")
	for _, commit := range commits {
		fmt.Println("%s", commit)
	}*/

	if !Helpers.IsDirExists(path) {
		err := Helpers.CreateNewDir(path)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	err := Cmd.Clone(path, resource)
	if err != nil {
		fmt.Println(err.Error())
	}

	//go Helpers.Exec(fmt.Sprintf("git clone %s %s", resource, path))

	time.Sleep(5 * time.Second)
}

