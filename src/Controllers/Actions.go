package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Shitovdm/GitRsync/src/components/cmd"
	"github.com/Shitovdm/GitRsync/src/components/configuration"
	"github.com/Shitovdm/GitRsync/src/components/helpers"
	"github.com/Shitovdm/GitRsync/src/components/logger"
	"github.com/Shitovdm/GitRsync/src/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ActionsController struct{}

func (ctrl ActionsController) Pull(c *gin.Context) {

	var pullActionRequest models.PullActionRequest
	conn, err := helpers.WsHandler(c.Writer, c.Request, &pullActionRequest)
	if err != nil {
		UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, StatusPullFailed)
		logger.Error("ActionsController/Pull", err.Error())
		logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, StatusPendingPull)
	logger.Trace("ActionsController/Pull", fmt.Sprintf("Start pulling repository with UUID %s...", pullActionRequest.RepositoryUuid))

	repositoryConfig := configuration.GetRepositoryByUuid(pullActionRequest.RepositoryUuid)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, StatusPullFailed)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", pullActionRequest.RepositoryUuid)
		logger.Error("ActionsController/Pull", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	platformConfig := configuration.GetPlatformByUuid(repositoryConfig.SourcePlatformUuid)
	if platformConfig == nil {
		UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, StatusPullFailed)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.SourcePlatformUuid)
		logger.Error("ActionsController/Pull", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	repositoryName := configuration.GetRepositorySourceRepositoryName(repositoryConfig)

	isNewRepository := false
	if !helpers.IsDirExists(configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name))) ||
		!helpers.IsDirExists(configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s\source\%s`, repositoryConfig.Name, repositoryName))) {
		logger.Trace("ActionsController/Pull", fmt.Sprintf("Repository %s has not been initialized earlier! Initialization...", repositoryConfig.Name))
		err = helpers.CreateNewDir(configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s\source`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, StatusPullFailed)
			ErrorMsg := fmt.Sprintf(`Error while creating new folder .\projects\%s\source`, repositoryConfig.Name)
			logger.Error("ActionsController/Pull", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
			logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
			return
		}
		logger.Trace("ActionsController/Pull", fmt.Sprintf("Root folders for repository %s succesfully created!", repositoryConfig.Name))
		isNewRepository = true
	}

	repositoryFullURL := platformConfig.Address + repositoryConfig.SourcePlatformPath
	repositoryFullPath := configuration.BuildPlatformPath(fmt.Sprintf(`projects\%s\source`, repositoryConfig.Name))
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
		UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, StatusPullFailed)
		Msg := "Error occurred while fetching source repository!"
		logger.Error("ActionsController/Pull", Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(Msg)))
		logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, StatusPulled)
	Msg := "Source repository fetched successfully!"
	logger.Success("ActionsController/Pull", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	return
}

func (ctrl ActionsController) Push(c *gin.Context) {

	var pushActionRequest models.PushActionRequest
	conn, err := helpers.WsHandler(c.Writer, c.Request, &pushActionRequest)
	if err != nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusPushFailed)
		logger.Error("ActionsController/Push", err.Error())
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusPendingPush)
	logger.Trace("ActionsController/Push", fmt.Sprintf("Start pushing repository with UUID %s...", pushActionRequest.RepositoryUuid))

	repositoryConfig := configuration.GetRepositoryByUuid(pushActionRequest.RepositoryUuid)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusPushFailed)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", pushActionRequest.RepositoryUuid)
		logger.Error("ActionsController/Push", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	platformConfig := configuration.GetPlatformByUuid(repositoryConfig.DestinationPlatformUuid)
	if platformConfig == nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusPushFailed)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.DestinationPlatformUuid)
		logger.Error("ActionsController/Push", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	//	Step 1. Cloning/pulling destination repository.
	destinationRepositoryName := configuration.GetRepositoryDestinationRepositoryName(repositoryConfig)
	sourceRepositoryName := configuration.GetRepositorySourceRepositoryName(repositoryConfig)

	isNewRepository := false
	if !helpers.IsDirExists(configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name))) ||
		!helpers.IsDirExists(configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s\destination\%s`, repositoryConfig.Name, destinationRepositoryName))) {
		logger.Trace("ActionsController/Push", fmt.Sprintf("Repository %s has not been initialized earlier! Initializing...", repositoryConfig.Name))
		err = helpers.CreateNewDir(configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s\destination`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusPushFailed)
			ErrorMsg := fmt.Sprintf("Error while creating new folder ./projects/%s/destination", repositoryConfig.Name)
			logger.Error("ActionsController/Push", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
			logger.Warning("ActionsController/Push", "Remote repository update aborted!")
			return
		}
		logger.Trace("ActionsController/Push", fmt.Sprintf("Root folder for repository %s successfully created!", repositoryConfig.Name))
		isNewRepository = true
	}

	repositoryFullURL := platformConfig.Address + repositoryConfig.DestinationPlatformPath
	repositoryFullPath := configuration.BuildPlatformPath(fmt.Sprintf(`projects\%s`, repositoryConfig.Name))
	finishFetching := make(chan bool)
	if isNewRepository {
		//	Clone action.
		logger.Trace("ActionsController/Push", fmt.Sprintf("Cloning repository from %s...", repositoryFullURL))
		go func() {
			cloneResult := cmd.Clone(repositoryFullPath+`\destination\`, repositoryFullURL)
			if cloneResult {
				logger.Trace("ActionsController/Push", fmt.Sprintf("Repository %s cloned successfully!", repositoryFullURL))
				finishFetching <- true
			} else {
				logger.Error("ActionsController/Push", fmt.Sprintf("Error occurred while cloning repository %s!", repositoryFullURL))
				finishFetching <- false
			}
		}()
	} else {
		//	Pull action.
		logger.Trace("ActionsController/Push", fmt.Sprintf("Pulling new from %s...", repositoryFullURL))
		go func() {
			pullResult := cmd.Pull(repositoryFullPath + `\destination\` + destinationRepositoryName)
			if pullResult {
				logger.Trace("ActionsController/Push", fmt.Sprintf("Repository %s pulled successfully!", repositoryFullURL))
				finishFetching <- true
			} else {
				logger.Error("ActionsController/Push", fmt.Sprintf("Error occurred while pulling repository %s!", repositoryFullURL))
				finishFetching <- false
			}
		}()
	}

	fetchRes := <-finishFetching
	if !fetchRes {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusPushFailed)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError("Error occurred while fetching destination repository!")))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	//	Step 2. Copying all files from source to destination repositories dir`s.
	logger.Trace("ActionsController/Push", "Start copying all repository files...")
	if !cmd.CopyRepository(repositoryFullPath, destinationRepositoryName, sourceRepositoryName) {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusPushFailed)
		logger.Error("ActionsController/Push", "Error occurred while copying repository files!")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError("Error occurred while copying repository files!")))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}
	logger.Trace("ActionsController/Push", "Repository files successfully copied!")

	destinationRepositoryPath := repositoryFullPath + `\destination\` + destinationRepositoryName

	//	Step 3. Checking needed pushing destination repository.
	commits, err := cmd.Log(destinationRepositoryPath, "origin/master..HEAD", -1)
	if err == nil {
		if len(commits) == 0 {
			UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusSynchronized)
			logger.Debug("ActionsController/Push", "Destination repository does not need to be updated, all changes are pushed earlier!")
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess("Destination repository does not need to be updated!", nil)))
			return
		}
	}
	logger.Trace("ActionsController/Push", "Remote destination repository needs updating, have unpushed changes!")

	//	Step 4. Rewriting commits author (if needed).
	appConfig := configuration.GetAppConfig()
	if appConfig.CommitsOverriding.State {
		logger.Trace("ActionsController/Pull", "Overriding source repository commits author...")
		if !cmd.OverrideAuthor(destinationRepositoryPath, appConfig.CommitsOverriding) {
			UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusPushFailed)
			Msg := "Error occurred while overriding destination repository commits author!"
			logger.Error("ActionsController/Pull", Msg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(Msg)))
			logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
			return
		}
		logger.Trace("ActionsController/Pull", "All commits in source repository successfully overridden!")
	} else {
		logger.Trace("ActionsController/Pull", "Overriding source repository commits not needed!")
	}

	//	Step 5. Pushing destination repository to remote.
	logger.Trace("ActionsController/Push", "Pushing destination repository...")
	if !cmd.Push(destinationRepositoryPath) {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusPushFailed)
		logger.Error("ActionsController/Push", "Error occurred while pushing destination repository!")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError("Error occurred while pushing destination repository!")))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, StatusSynchronized)
	Msg := "Destination repository successfully pushed!"
	logger.Success("ActionsController/Push", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	return
}

func (ctrl ActionsController) Clear(c *gin.Context) {

	var cleanActionRequest models.CleanActionRequest
	conn, err := helpers.WsHandler(c.Writer, c.Request, &cleanActionRequest)
	if err != nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, StatusCleanFailed)
		logger.Error("ActionsController/Clear", err.Error())
		logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
		return
	}

	UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, StatusPendingClean)
	logger.Trace("ActionsController/Clear", fmt.Sprintf("Start cleaning repository with UUID %s...", cleanActionRequest.RepositoryUuid))

	repositoryConfig := configuration.GetRepositoryByUuid(cleanActionRequest.RepositoryUuid)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, StatusCleanFailed)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", cleanActionRequest.RepositoryUuid)
		logger.Error("ActionsController/Clear", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
		return
	}

	platformConfig := configuration.GetPlatformByUuid(repositoryConfig.SourcePlatformUuid)
	if platformConfig == nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, StatusCleanFailed)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.SourcePlatformUuid)
		logger.Error("ActionsController/Clear", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
		return
	}

	if helpers.IsDirExists(configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name))) {
		err = helpers.RemoveDir(configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, StatusCleanFailed)
			ErrorMsg := "Error deleting project folder! Check if the folder is used by other programs and try again."
			logger.Error("ActionsController/Clear", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
			logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
			return
		}
	}

	UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, StatusCleaned)
	Msg := "Repository runtime data successfully cleaned!"
	logger.Success("ActionsController/Clear", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	return
}

func (ctrl ActionsController) Info(c *gin.Context) {
	var infoActionRequest models.InfoActionRequest
	conn, err := helpers.WsHandler(c.Writer, c.Request, &infoActionRequest)
	if err != nil {
		logger.Error("ActionsController/Info", err.Error())
		logger.Warning("ActionsController/Info", "Getting repository info aborted!")
		return
	}

	repositoryConfig := configuration.GetRepositoryByUuid(infoActionRequest.RepositoryUuid)
	if repositoryConfig == nil {
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", infoActionRequest.RepositoryUuid)
		logger.Error("ActionsController/Info", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		logger.Warning("ActionsController/Info", "Getting repository info aborted!")
		return
	}

	destinationRepositoryName := configuration.GetRepositoryDestinationRepositoryName(repositoryConfig)
	repositoryFullPath := configuration.BuildPlatformPath(fmt.Sprintf(`projects\%s`, repositoryConfig.Name))
	destinationRepositoryPath := repositoryFullPath + `\destination\` + destinationRepositoryName
	commitsLimit := configuration.GetAppConfigField("Common", "RecentCommitsShown")
	commits, err := cmd.Log(destinationRepositoryPath, "", int(commitsLimit.Int()))
	if err != nil {
		UpdateRepositoryStatus(infoActionRequest.RepositoryUuid, StatusFailed)
		ErrorMsg := fmt.Sprintf("Unable to select commits for repository %s!", destinationRepositoryName)
		logger.Error("ActionsController/Info", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		logger.Warning("ActionsController/Info", "Getting repository info aborted!")
		return
	}

	commitsJson, _ := json.Marshal(commits)
	Msg := "Repository commits list selected successfully!"
	logger.Success("ActionsController/Info", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, string(commitsJson))))
	return
}

func (ctrl ActionsController) Block(c *gin.Context) {
	var blockActionRequest models.BlockActionRequest
	conn, err := helpers.WsHandler(c.Writer, c.Request, &blockActionRequest)
	if err != nil {
		logger.Error("ActionsController/Block", err.Error())
		logger.Warning("ActionsController/Block", "Repository blocking aborted!")
		return
	}

	UpdateRepositoryState(blockActionRequest.RepositoryUuid, models.StateBlocked)
	Msg := "Repository successfully blocked!"
	logger.Success("ActionsController/Block", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	return
}

func (ctrl ActionsController) Activate(c *gin.Context) {
	var activateActionRequest models.ActivateActionRequest
	conn, err := helpers.WsHandler(c.Writer, c.Request, &activateActionRequest)
	if err != nil {
		logger.Error("ActionsController/Activate", err.Error())
		logger.Warning("ActionsController/Activate", "Repository activate aborted!")
		return
	}

	UpdateRepositoryState(activateActionRequest.RepositoryUuid, models.StateActive)
	Msg := "Repository successfully activated!"
	logger.Success("ActionsController/Activate", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	return
}

func BuildWsJsonError(message string) string {
	return fmt.Sprintf(`{"status":"error","message":"%s"}`, message)
}

func BuildWsJsonSuccess(message string, data interface{}) string {
	if data != nil {
		return fmt.Sprintf(`{"status":"success","message":"%s","data":%v}`, message, data)
	}
	return fmt.Sprintf(`{"status":"success","message":"%s","data":""}`, message)
}
