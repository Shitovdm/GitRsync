package actions

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/cmd"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/model/repository"
)

// SyncTags syncs repositories tags.
func SyncTags(repo *repository.Repository) error {

	repositoryFullPath := conf.BuildPlatformPath(fmt.Sprintf(`projects\%s`, repo.GetName()))
	sourceRepositoryName := repo.GetSourceRepositoryName()
	destinationRepositoryName := repo.GetDestinationRepositoryName()
	sourceRepositoryPath := repositoryFullPath + `\source\` + sourceRepositoryName
	destinationRepositoryPath := repositoryFullPath + `\destination\` + destinationRepositoryName

	sourceTags, err := cmd.GetTags(sourceRepositoryPath, -1)
	if err != nil {
		return fmt.Errorf("Unable to select tags for source repository %s! %s ",
			destinationRepositoryName, err.Error())
	}

	sourceCommits, err := cmd.Log(sourceRepositoryPath, "", -1)
	if err != nil {
		return fmt.Errorf("Unable to select commits for source repository %s! %s ", sourceRepositoryName, err.Error())
	}

	destinationCommits, err := cmd.Log(destinationRepositoryPath, "", -1)
	if err != nil {
		return fmt.Errorf("Unable to select commits for destination repository %s! %s ",
			destinationRepositoryName, err.Error())
	}

	sourceTags = cmd.ConvertTagsMeta(sourceCommits, destinationCommits, sourceTags)
	err = cmd.MakeTags(`C:\Users\Дмитрий\AppData\Roaming\GitRsync\projects\serv-queue-proxy\destination\serv-queue-proxy`, sourceTags)
	if err != nil {
		return fmt.Errorf("Unable to make new tags for destination repository %s! %s ",
			destinationRepositoryName, err.Error())
	}

	return nil
}
