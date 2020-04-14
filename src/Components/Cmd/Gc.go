package Cmd

import "os/exec"

func Gc() {
	exec.Command("git", "gc")
	return
}