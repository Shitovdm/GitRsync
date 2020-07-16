package actions

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/cmd"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/Shitovdm/GitRsync/src/model/platform"
	"github.com/Shitovdm/GitRsync/src/model/repository"
)

// PullCommits pulls source repository.
func PullCommits(repo *repository.Repository, re string) error {

	logger.Trace("ActionsController/Pull", fmt.Sprintf("Start pulling repository with UUID %s...", repo.UUID))
	repo.SetStatus(repository.StatusPendingPull)
	_ = repo.Update()
	sourceName := repo.GetSourceRepositoryName()

	isNewRepository := false
	if !helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repo.GetName()))) ||
		!helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s\source\%s`, repo.GetName(), sourceName))) {
		logger.Trace("ActionsController/Pull", fmt.Sprintf("Repository %s has not been initialized earlier! Initialization...", repo.GetName()))
		err := helper.CreateNewDir(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s\source`, repo.GetName())))
		if err != nil {
			repo.SetStatus(repository.StatusPullFailed)
			_ = repo.Update()
			return fmt.Errorf(`Error while creating new folder .\projects\%s\source `, repo.GetName())
		}
		logger.Trace("ActionsController/Pull", fmt.Sprintf("Root folders for repository %s succesfully created!", repo.GetName()))
		isNewRepository = true
	}

	pl := platform.Get(repo.DestinationPlatformUUID)
	repositoryFullURL := pl.GetAddress() + repo.GetDestinationPlatformPath()
	repositoryFullPath := conf.BuildPlatformPath(fmt.Sprintf(`projects\%s\destination`, repo.GetName()))
	if re == "source" {
		pl = platform.Get(repo.SourcePlatformUUID)
		repositoryFullURL = pl.GetAddress() + repo.GetSourcePlatformPath()
		repositoryFullPath = conf.BuildPlatformPath(fmt.Sprintf(`projects\%s\source`, repo.GetName()))
	}

	finishFetching := make(chan bool)
	if isNewRepository {
		//	Clone action.
		logger.Trace("ActionsController/Pull", fmt.Sprintf("Cloning repository from %s...", repositoryFullURL))
		go func() {
			cloneResult := cmd.Clone(repositoryFullPath, repositoryFullURL)
			if cloneResult {
				logger.Trace("ActionsController/Pull", fmt.Sprintf("Repository %s cloned successfully!", repositoryFullURL))
				finishFetching <- true
			} else {
				logger.Error("ActionsController/Pull", fmt.Sprintf("Error occurred while cloning repository %s!", repositoryFullURL))
				finishFetching <- false
			}
		}()
	} else {
		//	Pull action.
		logger.Trace("ActionsController/Pull", fmt.Sprintf("Fetching new from %s...", repositoryFullURL))
		go func() {
			pullResult := cmd.Pull(repositoryFullPath + `\` + repo.GetName())
			if pullResult {
				logger.Trace("ActionsController/Pull", fmt.Sprintf("Repository %s pulled successfully!", repositoryFullURL))
				finishFetching <- true
			} else {
				logger.Error("ActionsController/Pull", fmt.Sprintf("Error occurred while pulling repository %s!", repositoryFullURL))
				finishFetching <- false
			}
		}()
	}

	fetchRes := <-finishFetching
	if !fetchRes {
		repo.SetStatus(repository.StatusPullFailed)
		_ = repo.Update()
		return fmt.Errorf(`Error occurred while fetching source repository %s! `, repo.GetName())
	}

	return nil
}
