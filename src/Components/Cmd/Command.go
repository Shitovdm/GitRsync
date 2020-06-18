package cmd

import "os/exec"

func command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	return cmd
}
