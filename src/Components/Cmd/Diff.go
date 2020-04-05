package Cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)
func Diff(path string) bool {

	if path == "" {
		return false
	}

	var cmd *exec.Cmd
	cmd = command("git", "diff")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	lineBytes := bytes.Split(output, []byte{'\n'})
	lineBytes = lineBytes[0 : len(lineBytes)-1]

	for x := 0; x < len(lineBytes); x++ {
		fmt.Println(lineBytes[x])
		/*if string(lineBytes[x]) == "nothing to commit, working tree clean" {
			return true
		}*/
	}

	return false
}