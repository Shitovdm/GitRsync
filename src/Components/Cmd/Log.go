package Cmd

import (
	"bytes"
	"os/exec"
)

func Log(path string, hashish string) ([]*Commit, error) {

	var cmd *exec.Cmd
	if hashish == "" {
		cmd = command("git", "log", "--pretty=%H")
	} else {
		cmd = command("git", "log", "--pretty=%H", hashish)
	}
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return []*Commit{}, err
	}

	lineBytes := bytes.Split(output, []byte{'\n'})
	lineBytes = lineBytes[0 : len(lineBytes)-1]
	commits := make([]*Commit, len(lineBytes))

	for x := 0; x < len(lineBytes); x++ {
		commit, commitErr := NewCommit(path, string(lineBytes[x]))
		if commitErr != nil {
			return []*Commit{}, commitErr
		}
		commits[x] = commit
	}

	return commits, nil
}
