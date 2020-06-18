package cmd

import "os/exec"

// Gc runs git GC.
func Gc() {
	exec.Command("git", "gc")
}
