package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Tag contains tag metadata.
type Tag struct {
	Ref         string
	Hash        string
	Author      string
	AuthorEmail string
	Subject     string
	Time        string
	CommitHash  string
}

// String returns tag metadata as string.
func (tag *Tag) String() string {
	return fmt.Sprintf("\n+ Tag: %s\n| Hash: %s\n| Author: %s %s\n| Time: %s\n| Subject: %s\n| CommitHash: %s",
		tag.Ref, tag.Hash, tag.Author, tag.AuthorEmail, tag.Time, tag.Subject, tag.CommitHash)
}

// GetTags returns all tags metadata in []Tag struct.
func GetTags(path string, limit int) ([]*Tag, error) {

	logFormat := "%(refname)<|>%(*objectname)<|>%(*authorname)<|>%(*authoremail)<|>%(*subject)<|>%(*authordate)<|>"
	cmd := command("git", "tag", fmt.Sprintf("--format=\"%s\"", logFormat))
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return []*Tag{}, err
	}

	lineBytes := bytes.Split(output, []byte{'\n'})
	lineBytes = lineBytes[0 : len(lineBytes)-1]

	tags := make([]*Tag, len(lineBytes))
	if limit != -1 {
		tags = make([]*Tag, limit)
	}

	for x := 0; x < len(lineBytes); x++ {
		if limit != -1 && x >= limit {
			break
		}

		tag := &Tag{}
		rawMeta := strings.Split(string(lineBytes[x]), "<|>")
		tag.Ref = strings.Split(rawMeta[0], "/")[2]
		tag.Hash = rawMeta[1]
		tag.Author = rawMeta[2]
		tag.AuthorEmail = rawMeta[3]
		tag.Subject = rawMeta[4]
		tag.Time = rawMeta[5]

		commitHash, err := GetTagCommitHash(path, tag)
		if err != nil {
			return []*Tag{}, err
		}
		tag.CommitHash = *commitHash
		tags[x] = tag
	}

	return tags, nil
}

// GetTagCommitHash returns tag ref commit hash.
func GetTagCommitHash(path string, tag *Tag) (*string, error) {

	cmd := command("git", "show", tag.Ref)
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	rawMeta := string(output)
	start := strings.Index(rawMeta, "commit ")
	end := strings.Index(rawMeta, "Author:")
	if start == -1 || start == 0 || end == -1 || end == 0 {
		return nil, errors.New("Out of range! ")
	}
	searchMatch := rawMeta[(start + 7):(end - 1)]

	return &searchMatch, nil
}

// MakeTag creates project tag.
func MakeTag(path string, tag *Tag) error {

	fmt.Println("Making tag:", tag)
	cmd := command("git", "tag", "-a", tag.Ref, "-m", "'"+tag.Subject+"'", tag.CommitHash)
	cmd.Dir = path
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

// MakeTags creates project tags.
func MakeTags(path string, tags []*Tag) error {

	for _, tag := range tags {
		err := MakeTag(path, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

// ConvertTagsMeta re-bulds tags array. Places commit hashes for destination repository.
func ConvertTagsMeta(sourceCommits, destinationCommits []*Commit, tags []*Tag) []*Tag {

	for i, tag := range tags {
		dstSubject := ""
		dstTime := 0
		for _, sourceCommit := range sourceCommits {
			if sourceCommit.Hash == tag.CommitHash {
				dstSubject = sourceCommit.Subject
				dstTime = sourceCommit.Timestamp
				break
			}
		}
		for _, destinationCommit := range destinationCommits {
			if destinationCommit.Timestamp == dstTime && destinationCommit.Subject == dstSubject {
				fmt.Println("change tag", tag.Ref, "commit hash", tag.CommitHash, "to", destinationCommit.Hash)
				tags[i].CommitHash = destinationCommit.Hash
				break
			}
		}
	}

	return tags
}
