package Cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func Clone(path string, url string) error {

	if url == "" {
		return errors.New("Passed invalid URL! ")
	}

	var cmd *exec.Cmd
	cmd = command("git", "clone", url)
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lineBytes := bytes.Split(output, []byte{'\n'})
	fmt.Println(lineBytes)
	// The last split is just an empty string, right?
	/*lineBytes = lineBytes[0 : len(lineBytes)-1]
	commits := make([]*Commit, len(lineBytes))

	for x := 0; x < len(lineBytes); x++ {
		commit, commitErr := NewCommit(path, string(lineBytes[x]))
		if commitErr != nil {
			return commitErr
		}
		commits[x] = commit
	}*/

	return nil
}
