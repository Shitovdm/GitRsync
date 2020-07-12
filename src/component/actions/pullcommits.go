package actions

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/cmd"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/Shitovdm/GitRsync/src/model/repository"
	"github.com/gorilla/websocket"
)

func PullCommits(repositoryUUID string) error {

	repo := repository.Get(repositoryUUID)
	repo.SetStatus(repository.StatusPendingPull)
	_ = repo.Update()

	logger.Trace("ActionsController/Pull", fmt.Sprintf("Start pulling repository with UUID %s...", repositoryUUID))

	platformConfig := conf.GetPlatformByUUID(repositoryConfig.SourcePlatformUUID)
	if platformConfig == nil {
		UpdateRepositoryStatus(repositoryUUID, repository.StatusPullFailed)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.SourcePlatformUUID)
		logger.Error("ActionsController/Pull", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
		logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	sourceName := repo.GetSourceRepositoryName()
	repositoryName := conf.GetRepositorySourceRepositoryName(repositoryConfig)

	isNewRepository := false
	if !helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repo.GetName()))) ||
		!helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s\source\%s`, repo.GetName(), repositoryName))) {
		logger.Trace("ActionsController/Pull", fmt.Sprintf("Repository %s has not been initialized earlier! Initialization...", repositoryConfig.Name))
		err = helper.CreateNewDir(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s\source`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(repositoryUUID, repository.StatusPullFailed)
			ErrorMsg := fmt.Sprintf(`Error while creating new folder .\projects\%s\source`, repositoryConfig.Name)
			logger.Error("ActionsController/Pull", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
			logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
			return
		}
		logger.Trace("ActionsController/Pull", fmt.Sprintf("Root folders for repository %s succesfully created!", repositoryConfig.Name))
		isNewRepository = true
	}

	repositoryFullURL := platformConfig.Address + repositoryConfig.SourcePlatformPath
	repositoryFullPath := conf.BuildPlatformPath(fmt.Sprintf(`projects\%s\source`, repositoryConfig.Name))
	//
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
			pullResult := cmd.Pull(repositoryFullPath + `\` + repositoryName)
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
		UpdateRepositoryStatus(pullActionRequest.RepositoryUUID, repository.StatusPullFailed)
		Msg := "Error occurred while fetching source repository!"
		logger.Error("ActionsController/Pull", Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(Msg)))
		logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	return nil
}
