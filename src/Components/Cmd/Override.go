package Cmd

import (
	"fmt"
	"os/exec"
	"time"
)

func Override(path string, username string, email string) bool {



	
	fmt.Println("Status: ", Status(path))

	commits, err := Log(path, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(commits)

	var cmd *exec.Cmd
	//author := fmt.Sprintf("GIT_AUTHOR_NAME='%s';GIT_AUTHOR_EMAIL='%s'", username, email)
	//committer := fmt.Sprintf("GIT_COMMITTER_NAME='%s';GIT_COMMITTER_EMAIL='%s'", username, email)
	//envFilter := fmt.Sprintf("\"%s;%s;\"", author, committer)

	gitCmd := fmt.Sprintf(`git filter-branch --env-filter '
		GIT_AUTHOR_NAME="%s";
		GIT_AUTHOR_EMAIL="%s";
		GIT_COMMITTER_NAME="%s";
		GIT_COMMITTER_EMAIL="%s";' HEAD;`,
		username,
		email,
		username,
		email)

	cmd = exec.Command("bash", "-c", gitCmd)

	//cmd = command("git", "filter-branch", "--force", "[--env-filter "+envFilter+"]")
	//cmd = exec.Command("git", "filter-branch", "--force", "[--env-filter "+fmt.Sprintf("\"GIT_AUTHOR_NAME='%s'\"]", "Shitov Dmitry"))
	//cmd = command("git", "filter-branch", "--env-filter \"GIT_AUTHOR_NAME='Shitov Dmitry';GIT_AUTHOR_EMAIL='shitov.dm@gmail.com';GIT_COMMITTER_NAME='Shitov Dmitry';GIT_COMMITTER_EMAIL='shitov.dm@gmail.com';\"")
	cmd.Dir = path
	StdoutPipe, err := cmd.StderrPipe()
	if err != nil {
		return false
	}

	fmt.Println("DIR: " + cmd.Dir)
	fmt.Println("PATH: " + cmd.Path)
	fmt.Println("CMD command: " + cmd.String())

	breakFlag := false
	finish := make(chan bool)
	go func() {
		go func() {
			for {
				if breakFlag {
					break
				}
				output := make([]byte, 256, 256)
				_, _ = StdoutPipe.Read(output)
				fmt.Println(string(output))
				if string(output) == "Ref 'refs/heads/master' was rewritten" ||
					string(output) == "WARNING: Ref 'refs/heads/master' is unchanged" {
					finish <- true
				}
				if string(output) == "exit status 128" {
					finish <- false
				}

				time.Sleep(50 * time.Millisecond)
			}
		}()

		err = cmd.Run()
		if err != nil {
			time.Sleep(5*time.Second)
			fmt.Println("running error!" + err.Error())
			breakFlag = true
			finish <- false
		}

		_ = cmd.Wait()

		finish <- true
	}()

	result := <-finish
	breakFlag = true
	return result
}
