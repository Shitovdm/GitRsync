package helper

import (
	"fmt"
	"github.com/speedata/gogit"
	"log"
	"os/exec"
	"regexp"
)

// Repo contain repository object.
type Repo struct {
	repository *gogit.Repository
}

// isEqual compares commits.
func isEqual(c1, c2 *gogit.Commit) bool {
	return c1.Oid == c2.Oid &&
		c1.Author.Name == c2.Author.Name &&
		c1.Author.Email == c2.Author.Email &&
		c1.Author.When == c2.Author.When &&
		c1.Committer.Name == c2.Committer.Name &&
		c1.Committer.Email == c2.Committer.Email &&
		c1.Committer.When == c2.Committer.When
}

// getCommitChainLength returns commit chain len.
func getCommitChainLength(ci *gogit.Commit, n int) int {
	for i := 1; i < n; i++ {
		ci = ci.Parent(0)
		if ci.ParentCount() == 0 {
			return i + 1
		}
	}

	return n
}

// GetLog returns log.
func (r *Repo) GetLog(n int) ([]*gogit.Commit, error) {
	ref, err := r.repository.LookupReference("HEAD")
	if err != nil {
		return nil, err
	}
	ci, err := r.repository.LookupCommit(ref.Oid)
	if err != nil {
		return nil, err
	}

	n = getCommitChainLength(ci, n)

	commitList := make([]*gogit.Commit, n)
	commitList[0] = ci

	for i := 1; i < n; i++ {
		ci = ci.Parent(0)
		if ci == nil {
			break
		}
		commitList[i] = ci
	}

	return commitList, nil
}

// IsDirty returns diff "dirty" flag.
func (r *Repo) IsDirty() bool {
	gitCmd := `[[ $(git diff --shortstat 2> /dev/null | tail -n1) != "" ]] && echo "dirty"`
	cmd := exec.Command("bash", "-c", gitCmd)
	output, _ := cmd.Output()

	re := regexp.MustCompile("dirty")
	return re.Match(output)
}

// SaveCommitIfModified returns the ref of change and error.
func (r *Repo) SaveCommitIfModified(commit *gogit.Commit) (string, error) {
	original, err := r.repository.LookupCommit(commit.Oid)
	if err != nil {
		return "", fmt.Errorf("Error finding matching commit: %s ", err)
	}

	if !isEqual(commit, original) {
		return r.SaveCommit(commit)
	}

	return "", nil
}

// SaveCommit save commit meta.
func (r *Repo) SaveCommit(commit *gogit.Commit) (string, error) {
	scope := ""
	if commit.Parent(0) != nil {
		scope = fmt.Sprintf("%s..HEAD", commit.Parent(0).Oid.String())
	}

	gitCmd := fmt.Sprintf(`git filter-branch --env-filter 'if [ $GIT_COMMIT = %s ]
		then
			export GIT_AUTHOR_NAME="%s" &&
			export GIT_AUTHOR_EMAIL="%s" &&
			export GIT_AUTHOR_DATE="%s" &&
			export GIT_COMMITTER_NAME="%s" &&
			export GIT_COMMITTER_EMAIL="%s" &&
			export GIT_COMMITTER_DATE="%s"; fi' %s &&
		rm -fr "$(git rev-parse --git-dir)/refs/original/"`,
		commit.Oid.String(),
		commit.Author.Name,
		commit.Author.Email,
		commit.Author.When.String(),
		commit.Committer.Name,
		commit.Committer.Email,
		commit.Committer.When.String(),
		scope)
	log.Println(gitCmd)
	cmd := exec.Command("bash", "-c", gitCmd)

	output, err := cmd.Output()
	re := regexp.MustCompile("Ref '([^']*)' was rewritten")
	match := re.FindStringSubmatch(string(output))

	if err != nil {
		return "", err
	}

	if len(match) == 0 {
		return "", fmt.Errorf("Git rewrite failed due to no change ")
	}

	return match[1], nil
}
