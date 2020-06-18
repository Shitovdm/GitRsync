package cmd

import "os/exec"

// command executes shell command.
func command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	return cmd
}
