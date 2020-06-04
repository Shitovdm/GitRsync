package Cmd

import (
	"fmt"
	"os/exec"
	"time"
)

func OverrideAuthor(path string, username string, email string) bool {

	var cmd *exec.Cmd

	cmd = exec.Command("bash", "-c", "git status")
	cmd.Dir = path
	_, err := cmd.Output()
	if err != nil {
		return false
	}

	cmd = exec.Command("bash", "-c", "git diff")
	cmd.Dir = path
	_, err = cmd.Output()
	if err != nil {
		return false
	}

	gitCmd := fmt.Sprintf(
		`git filter-branch -f --env-filter "GIT_AUTHOR_NAME='%s'; GIT_AUTHOR_EMAIL='%s'; GIT_COMMITTER_NAME='%s'; GIT_COMMITTER_EMAIL='%s';" HEAD;`,
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
