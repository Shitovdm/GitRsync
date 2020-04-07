package Cmd

import (
	"fmt"
	"os/exec"
	"time"
)

func Override(path string, username string, email string) bool {
	var cmd *exec.Cmd
	envFilter := fmt.Sprintf("'GIT_AUTHOR_NAME='%s';GIT_AUTHOR_EMAIL='%s';GIT_COMMITTER_NAME='%s';GIT_COMMITTER_EMAIL='%s';'", username, email, username, email)
	cmd = command("git", "filter-branch", "-f", "--env-filter", envFilter)
	cmd.Dir = path
	StdoutPipe, err := cmd.StderrPipe()
	if err != nil {
		return false
	}

	fmt.Println("CMD command: " + cmd.String())

	breakFlag := false
	finish := make(chan bool)
	go func() {
		go func() {
			for {
				if breakFlag {
					break
				}
				output := make([]byte, 128, 128)
				_, _ = StdoutPipe.Read(output)
				if string(output) == "Ref 'refs/heads/master' was rewritten" ||
					string(output) == "WARNING: Ref 'refs/heads/master' is unchanged" {
					finish <- true
				}
				if string(output) == "exit status 128" {
					finish <- false
				}
				fmt.Println(string(output))
				time.Sleep(50 * time.Millisecond)
			}
		}()

		err = cmd.Run()
		if err != nil {
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
