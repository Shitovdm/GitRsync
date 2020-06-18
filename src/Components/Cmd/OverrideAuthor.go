package Cmd

import (
	"fmt"
	"github.com/Shitovdm/git-rsync/src/Models"
	"os/exec"
	"time"
)

func OverrideAuthor(path string, CommitsOverridingConfig Models.CommitsOverriding) bool {

	var cmd *exec.Cmd

	cmd = exec.Command("bash", "-c", "git status")
	cmd.Dir = path
	_, err := cmd.Output()
	if err != nil {
		return false
	}

	cmd = exec.Command("bash", "-c", "git diff")
	cmd.Dir = path
	_, err = cmd.Output()
	if err != nil {
		return false
	}

	gitCmd := ``
	if CommitsOverridingConfig.OverrideCommitsWithOneAuthor {
		username := CommitsOverridingConfig.MasterUser.Username
		email := CommitsOverridingConfig.MasterUser.Email
		gitCmd = fmt.Sprintf(
			`git filter-branch -f --env-filter "GIT_AUTHOR_NAME='%s'; GIT_AUTHOR_EMAIL='%s'; GIT_COMMITTER_NAME='%s'; GIT_COMMITTER_EMAIL='%s';" HEAD;`,
			username, email, username, email)
	} else {
		gitCmd = fmt.Sprintf("`%s`", BuildFilterBranchExpression(CommitsOverridingConfig.CommittersRules))
	}

	cmd = exec.Command("bash", "-c", gitCmd)
	cmd.Dir = path
	StdoutPipe, err := cmd.StderrPipe()
	if err != nil {
		return false
	}

	breakFlag := false
	finish := make(chan bool)
	go func() {
		go func() {
			for {
				if breakFlag {
					break
				}
				output := make([]byte, 256, 256)
				_, _ = StdoutPipe.Read(output)
				raw := string(output)
				if raw == "Ref 'refs/heads/master' was rewritten" ||
					raw == "WARNING: Ref 'refs/heads/master' is unchanged" ||
					raw == "exit status 0" ||
					raw == "exit status 2" {
					finish <- true
				}
				if raw == "exit status 128" ||
					raw == "exit status 1" {
					//finish <- false
					finish <- true
				}

				time.Sleep(50 * time.Millisecond)
			}
		}()

		err = cmd.Run()
		if err != nil {
			fmt.Println("running error!" + err.Error())
			breakFlag = true
			//finish <- false
			finish <- true
		}

		_ = cmd.Wait()
		finish <- true
	}()

	result := <-finish
	breakFlag = true

	return result
}

func BuildFilterBranchExpression(committersRules []Models.CommittersRule) string {
	resultExpression := `git filter-branch --env-filter '`
	for _, rule := range committersRules {
		resultExpression += fmt.Sprintf(`
if test "$GIT_AUTHOR_NAME" = "%s"
then
	GIT_AUTHOR_NAME="%s"
fi
if test "$GIT_AUTHOR_EMAIL" = "%s"
then
	GIT_AUTHOR_EMAIL="%s"
fi
if test "$GIT_COMMITTER_NAME" = "%s"
then
	GIT_COMMITTER_NAME="%s"
fi
if test "$GIT_COMMITTER_EMAIL" = "%s"
then
	GIT_COMMITTER_EMAIL="%s"
fi`, rule.Old.Username, rule.New.Username, rule.Old.Email, rule.New.Email, rule.Old.Username, rule.New.Username, rule.Old.Email, rule.New.Email)
	}
	resultExpression += `' HEAD;`

	return resultExpression
}
