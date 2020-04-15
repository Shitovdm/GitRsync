package Cmd

import (
	"fmt"
	"os/exec"
	"time"
)

func ChangeCommit(path string, username string, email string) bool {
	var cmd *exec.Cmd

	cmd = exec.Command("bash", "-c", "git status")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		fmt.Println(fmt.Sprintf("%s", out))
	}

	cmd = exec.Command("bash", "-c", "git diff")
	cmd.Dir = path
	out, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		fmt.Println(fmt.Sprintf("%s", out))
	}

	gitCmd := fmt.Sprintf(
		`git filter-branch --env-filter "GIT_AUTHOR_NAME='%s'; GIT_AUTHOR_EMAIL='%s'; GIT_COMMITTER_NAME='%s'; GIT_COMMITTER_EMAIL='%s';"`,
		username, email, username, email)
	cmd = exec.Command("bash", "-c", gitCmd)
	cmd.Dir = path
	StdoutPipe, err := cmd.StderrPipe()
	if err != nil {
		return false
	}

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
