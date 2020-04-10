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

	//resource := "https://github.com/Shitovdm/rpc.git"
	//path := Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s", "rpc-test-sync"))


	/*commits, _ := Cmd.Log(path, "ba3edcc592c66d40b18613ac044d0bcf277eb08b")
	for _, commit := range commits {
		fmt.Println("%s", commit)
	}*/

	/*if !Helpers.IsDirExists(path) {
		err := Helpers.CreateNewDir(path)
		if err != nil {
			fmt.Println(err.Error())
		}
	}*/

	/*result := Cmd.Clone(path, resource)
	if result {
		fmt.Println("Repository cloned successfully!")
	}else{
		fmt.Println("Repository cloning error!")
	}*/

	/*path = Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s", "rpc-test-sync/rpc"))
	result := Cmd.Pull(path)
	if result {
		fmt.Println("Repository pull successfully!")
	}else{
		fmt.Println("Repository pulling error!")
	}*/

	//go Helpers.Exec(fmt.Sprintf("git clone %s %s", resource, path))




	//Cmd.Override("C:/Users/Дмитрий/AppData/Roaming/GitRsync/projects/lib-go-amqp-first/destination/lib-go-amqp-firs", "Shitov Dmitry", "shitov.dm@gmail.com")

	Application.StartServer()

}

