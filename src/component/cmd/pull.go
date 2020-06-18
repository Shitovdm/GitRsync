package cmd

import (
	"os/exec"
	"time"
)

func Pull(path string) bool {
	var cmd *exec.Cmd
	cmd = command("git", "pull")
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
				output := make([]byte, 128, 128)
				_, _ = StdoutPipe.Read(output)
				if string(output) == "fatal: destination path 'rpc' already exists and is not an empty directory." ||
					string(output) == "exit status 128" {
					finish <- false
				}
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
