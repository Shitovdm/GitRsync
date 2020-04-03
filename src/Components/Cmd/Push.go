package Cmd

import (
	"os/exec"
	"time"
)

func Push(path string) bool {

	var cmd *exec.Cmd
	cmd = command("git", "push", "origin", "master:master")
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

				//if string(output) == "fatal: destination path 'rpc' already exists and is not an empty directory." ||
				//	string(output) == "exit status 128" {
				if string(output)[:5] == "error" || string(output)[:5] == "fatal" || string(output)[:4] == "exit" {
					finish <- false
				}
				time.Sleep(500 * time.Millisecond)
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
	return result
}
