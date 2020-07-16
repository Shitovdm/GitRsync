package actions

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/cmd"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/model/repository"
)

// PushCommits updates remote repository.
func PushCommits(repo *repository.Repository) error {

	repositoryFullPath := conf.BuildPlatformPath(fmt.Sprintf(`projects\%s\destination`, repo.GetName()))
	destinationRepositoryName := repo.GetDestinationRepositoryName()
	destinationRepositoryPath := repositoryFullPath + `\destination\` + destinationRepositoryName

	if !cmd.Push(destinationRepositoryPath) {
		repo.SetStatus(repository.StatusPushFailed)
		_ = repo.Update()
		return fmt.Errorf(`Error occurred while pushing destination repository %s! `, repo.GetName())
	}

	return nil
}
