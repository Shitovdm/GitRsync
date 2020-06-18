package cmd

import (
	"bytes"
)

// Status returns git status.
func Status(path string) bool {

	if path == "" {
		return false
	}

	cmd := command("git", "status")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	lineBytes := bytes.Split(output, []byte{'\n'})
	lineBytes = lineBytes[0 : len(lineBytes)-1]

	for x := 0; x < len(lineBytes); x++ {
		if string(lineBytes[x]) == "nothing to commit, working tree clean" {
			return true
		}
	}

	return false
}
