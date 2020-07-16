package actions

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/cmd"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/model/repository"
)

func GetCommits(repo *repository.Repository, limit int) ([]*cmd.Commit, error) {

	repositoryFullPath := conf.BuildPlatformPath(fmt.Sprintf(`projects\%s`, repo.GetName()))
	destinationRepositoryPath := repositoryFullPath + `\destination\` + repo.GetDestinationRepositoryName()
	commits, err := cmd.Log(destinationRepositoryPath, "", limit)
	if err != nil {
		return nil, err
	}

	return commits, nil
}
