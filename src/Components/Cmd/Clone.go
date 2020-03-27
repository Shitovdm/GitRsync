package Cmd

import (
	"os/exec"
	"time"
)

func Clone(path string, url string) bool {
	if url == "" {
		return false
	}

	var cmd *exec.Cmd
	cmd = command("git", "clone", url)
	cmd.Dir = path
	StdoutPipe, err := cmd.StderrPipe()
	if err != nil {
		return false
	}

	finish := make(chan bool)
	go func() {
		go func() {
			for {
				output := make([]byte, 128, 128)
				_, _ = StdoutPipe.Read(output)
				if string(output) == "fatal: destination path 'rpc' already exists and is not an empty directory." ||
					string(output) == "exit status 128" {
					finish <- false
				}
				time.Sleep(10 * time.Microsecond)
			}
		}()

		err = cmd.Run()
		if err != nil {
			finish <- false
		}

		_ = cmd.Wait()

		finish <- true
	}()

	result := <-finish
	return result
}
