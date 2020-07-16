package actions

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/cmd"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/model/repository"
)

func OverrideCommits(repo *repository.Repository) error {

	appConfig := conf.GetAppConfig()
	destinationRepositoryName := repo.GetDestinationRepositoryName()
	repositoryFullPath := conf.BuildPlatformPath(fmt.Sprintf(`projects\%s\destination`, repo.GetName()))
	destinationRepositoryPath := repositoryFullPath + `\destination\` + destinationRepositoryName
	if !cmd.OverrideAuthor(destinationRepositoryPath, appConfig.CommitsOverriding) {
		return fmt.Errorf("Error occurred while overriding destination repository commits author! ")
	}

	return nil
}
