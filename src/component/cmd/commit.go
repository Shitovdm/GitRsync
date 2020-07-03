package cmd

import (
	"bytes"
	"fmt"
	"strconv"
)

// Commit contains commit metadata.
type Commit struct {
	Hash        string
	Author      string
	AuthorEmail string
	ParentHash  string
	Subject     string
	Timestamp   int
}

// String returns commit metadata as string.
func (commit *Commit) String() string {
	return fmt.Sprintf("\n+ Commit: %s\n| Author: %s <%s>\n| Parent: %s\n"+
		"| Timestamp: %d\n| Subject: %s", commit.Hash, commit.Author,
		commit.AuthorEmail, commit.ParentHash, commit.Timestamp,
		commit.Subject)
}

// NewCommit returns commit metadata in Commit struct.
func NewCommit(path, hashish string) (*Commit, error) {
	logFormat := "%H%n%an%n%ae%n%ct%n%P%n%s%n%b"
	commit := &Commit{}

	cmd := command("git", "log", "-1", "--pretty="+logFormat, hashish)
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return &Commit{}, err
	}

	lineBytes := bytes.Split(output, []byte{'\n'})
	commit.Hash = string(bytes.TrimSpace(lineBytes[0]))
	commit.Author = string(bytes.TrimSpace(lineBytes[1]))
	commit.AuthorEmail = string(bytes.TrimSpace(lineBytes[2]))
	commit.Timestamp, err = strconv.Atoi(string(bytes.TrimSpace(lineBytes[3])))
	if err != nil {
		return &Commit{}, err
	}
	commit.ParentHash = string(bytes.TrimSpace(lineBytes[4]))
	commit.Subject = string(bytes.TrimSpace(lineBytes[5]))

	return commit, nil
}

// GetCommitHashByBodyAndTime commit hash by commit subject and time.
func GetCommitHashByBodyAndTime(path, subject string, timestamp int) (*string, error) {

	commits, err := Log(path, "origin/master..HEAD", -1)
	if err == nil {
		return nil, err
	}

	for _, commit := range commits {
		if commit.Subject == subject && commit.Timestamp == timestamp {
			return &commit.Hash, nil
		}
	}

	return nil, nil
}