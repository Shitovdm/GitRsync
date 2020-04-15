package Cmd

import (
	"fmt"
	"os/exec"
)

func ResetHard(path string) bool {

	if path == "" {
		return false
	}

	var cmd *exec.Cmd
	cmd = command("git", "reset --hard")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println(fmt.Sprintf("%s", out))
	return true
}