package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/cmd"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/Shitovdm/GitRsync/src/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ActionsController struct describes main actions section controller.
type ActionsController struct{}

// Pull describes pull repository action.
func (ctrl ActionsController) Pull(c *gin.Context) {

	var pullActionRequest model.PullActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &pullActionRequest)
	if err != nil {
		UpdateRepositoryStatus(pullActionRequest.RepositoryUUID, StatusPullFailed)
		logger.Error("ActionsController/Pull", err.Error())
		logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	UpdateRepositoryStatus(pullActionRequest.RepositoryUUID, StatusPendingPull)
	logger.Trace("ActionsController/Pull", fmt.Sprintf("Start pulling repository with UUID %s...", pullActionRequest.RepositoryUUID))

	repositoryConfig := conf.GetRepositoryByUUID(pullActionRequest.RepositoryUUID)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(pullActionRequest.RepositoryUUID, StatusPullFailed)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", pullActionRequest.RepositoryUUID)
		logger.Error("ActionsController/Pull", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
		logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	platformConfig := conf.GetPlatformByUUID(repositoryConfig.SourcePlatformUUID)
	if platformConfig == nil {
		UpdateRepositoryStatus(pullActionRequest.RepositoryUUID, StatusPullFailed)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.SourcePlatformUUID)
		logger.Error("ActionsController/Pull", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
		logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	repositoryName := conf.GetRepositorySourceRepositoryName(repositoryConfig)

	isNewRepository := false
	if !helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name))) ||
		!helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s\source\%s`, repositoryConfig.Name, repositoryName))) {
		logger.Trace("ActionsController/Pull", fmt.Sprintf("Repository %s has not been initialized earlier! Initialization...", repositoryConfig.Name))
		err = helper.CreateNewDir(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s\source`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(pullActionRequest.RepositoryUUID, StatusPullFailed)
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
		UpdateRepositoryStatus(pullActionRequest.RepositoryUUID, StatusPullFailed)
		Msg := "Error occurred while fetching source repository!"
		logger.Error("ActionsController/Pull", Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(Msg)))
		logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	UpdateRepositoryStatus(pullActionRequest.RepositoryUUID, StatusPulled)
	Msg := "Source repository fetched successfully!"
	logger.Success("ActionsController/Pull", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONSuccess(Msg, nil)))
}

// Push describes push repository action.
func (ctrl ActionsController) Push(c *gin.Context) {

	var pushActionRequest model.PushActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &pushActionRequest)
	if err != nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusPushFailed)
		logger.Error("ActionsController/Push", err.Error())
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusPendingPush)
	logger.Trace("ActionsController/Push", fmt.Sprintf("Start pushing repository with UUID %s...", pushActionRequest.RepositoryUUID))

	repositoryConfig := conf.GetRepositoryByUUID(pushActionRequest.RepositoryUUID)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusPushFailed)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", pushActionRequest.RepositoryUUID)
		logger.Error("ActionsController/Push", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	platformConfig := conf.GetPlatformByUUID(repositoryConfig.DestinationPlatformUUID)
	if platformConfig == nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusPushFailed)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.DestinationPlatformUUID)
		logger.Error("ActionsController/Push", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	//	Step 1. Cloning/pulling destination repository.
	destinationRepositoryName := conf.GetRepositoryDestinationRepositoryName(repositoryConfig)
	sourceRepositoryName := conf.GetRepositorySourceRepositoryName(repositoryConfig)

	isNewRepository := false
	if !helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name))) ||
		!helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s\destination\%s`, repositoryConfig.Name, destinationRepositoryName))) {
		logger.Trace("ActionsController/Push", fmt.Sprintf("Repository %s has not been initialized earlier! Initializing...", repositoryConfig.Name))
		err = helper.CreateNewDir(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s\destination`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusPushFailed)
			ErrorMsg := fmt.Sprintf("Error while creating new folder ./projects/%s/destination", repositoryConfig.Name)
			logger.Error("ActionsController/Push", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
			logger.Warning("ActionsController/Push", "Remote repository update aborted!")
			return
		}
		logger.Trace("ActionsController/Push", fmt.Sprintf("Root folder for repository %s successfully created!", repositoryConfig.Name))
		isNewRepository = true
	}

	repositoryFullURL := platformConfig.Address + repositoryConfig.DestinationPlatformPath
	repositoryFullPath := conf.BuildPlatformPath(fmt.Sprintf(`projects\%s`, repositoryConfig.Name))
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
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusPushFailed)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError("Error occurred while fetching destination repository!")))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	//	Step 2. Copying all files from source to destination repositories dir`s.
	logger.Trace("ActionsController/Push", "Start copying all repository files...")
	if !cmd.CopyRepository(repositoryFullPath, destinationRepositoryName, sourceRepositoryName) {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusPushFailed)
		logger.Error("ActionsController/Push", "Error occurred while copying repository files!")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError("Error occurred while copying repository files!")))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}
	logger.Trace("ActionsController/Push", "Repository files successfully copied!")

	destinationRepositoryPath := repositoryFullPath + `\destination\` + destinationRepositoryName

	//	Step 3. Checking needed pushing destination repository.
	commits, err := cmd.Log(destinationRepositoryPath, "origin/master..HEAD", -1)
	if err == nil {
		if len(commits) == 0 {
			UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusSynchronized)
			logger.Debug("ActionsController/Push", "Destination repository does not need to be updated, all changes are pushed earlier!")
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONSuccess("Destination repository does not need to be updated!", nil)))
			return
		}
	}
	logger.Trace("ActionsController/Push", "Remote destination repository needs updating, have unpushed changes!")

	//	Step 4. Rewriting commits author (if needed).
	appConfig := conf.GetAppConfig()
	if appConfig.CommitsOverriding.State {
		logger.Trace("ActionsController/Pull", "Overriding source repository commits author...")
		if !cmd.OverrideAuthor(destinationRepositoryPath, appConfig.CommitsOverriding) {
			UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusPushFailed)
			Msg := "Error occurred while overriding destination repository commits author!"
			logger.Error("ActionsController/Pull", Msg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(Msg)))
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
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusPushFailed)
		logger.Error("ActionsController/Push", "Error occurred while pushing destination repository!")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError("Error occurred while pushing destination repository!")))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, StatusSynchronized)
	Msg := "Destination repository successfully pushed!"
	logger.Success("ActionsController/Push", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONSuccess(Msg, nil)))
}

// Clear describes clear repository data action.
func (ctrl ActionsController) Clear(c *gin.Context) {

	var cleanActionRequest model.CleanActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &cleanActionRequest)
	if err != nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, StatusCleanFailed)
		logger.Error("ActionsController/Clear", err.Error())
		logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
		return
	}

	UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, StatusPendingClean)
	logger.Trace("ActionsController/Clear", fmt.Sprintf("Start cleaning repository with UUID %s...", cleanActionRequest.RepositoryUUID))

	repositoryConfig := conf.GetRepositoryByUUID(cleanActionRequest.RepositoryUUID)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, StatusCleanFailed)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", cleanActionRequest.RepositoryUUID)
		logger.Error("ActionsController/Clear", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
		logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
		return
	}

	platformConfig := conf.GetPlatformByUUID(repositoryConfig.SourcePlatformUUID)
	if platformConfig == nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, StatusCleanFailed)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.SourcePlatformUUID)
		logger.Error("ActionsController/Clear", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
		logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
		return
	}

	if helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name))) {
		err = helper.RemoveDir(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, StatusCleanFailed)
			ErrorMsg := "Error deleting project folder! Check if the folder is used by other programs and try again."
			logger.Error("ActionsController/Clear", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
			logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
			return
		}
	}

	UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, StatusCleaned)
	Msg := "Repository runtime data successfully cleaned!"
	logger.Success("ActionsController/Clear", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONSuccess(Msg, nil)))
}

// Info describes repository info action.
func (ctrl ActionsController) Info(c *gin.Context) {
	var infoActionRequest model.InfoActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &infoActionRequest)
	if err != nil {
		logger.Error("ActionsController/Info", err.Error())
		logger.Warning("ActionsController/Info", "Getting repository info aborted!")
		return
	}

	repositoryConfig := conf.GetRepositoryByUUID(infoActionRequest.RepositoryUUID)
	if repositoryConfig == nil {
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", infoActionRequest.RepositoryUUID)
		logger.Error("ActionsController/Info", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
		logger.Warning("ActionsController/Info", "Getting repository info aborted!")
		return
	}

	destinationRepositoryName := conf.GetRepositoryDestinationRepositoryName(repositoryConfig)
	repositoryFullPath := conf.BuildPlatformPath(fmt.Sprintf(`projects\%s`, repositoryConfig.Name))
	destinationRepositoryPath := repositoryFullPath + `\destination\` + destinationRepositoryName
	commitsLimit := conf.GetAppConfigField("Common", "RecentCommitsShown")
	commits, err := cmd.Log(destinationRepositoryPath, "", int(commitsLimit.Int()))
	if err != nil {
		UpdateRepositoryStatus(infoActionRequest.RepositoryUUID, StatusFailed)
		ErrorMsg := fmt.Sprintf("Unable to select commits for repository %s!", destinationRepositoryName)
		logger.Error("ActionsController/Info", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
		logger.Warning("ActionsController/Info", "Getting repository info aborted!")
		return
	}

	commitsJSON, _ := json.Marshal(commits)
	Msg := "Repository commits list selected successfully!"
	logger.Success("ActionsController/Info", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONSuccess(Msg, string(commitsJSON))))
}

// Block describes block repository action.
func (ctrl ActionsController) Block(c *gin.Context) {
	var blockActionRequest model.BlockActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &blockActionRequest)
	if err != nil {
		logger.Error("ActionsController/Block", err.Error())
		logger.Warning("ActionsController/Block", "Repository blocking aborted!")
		return
	}

	UpdateRepositoryState(blockActionRequest.RepositoryUUID, model.StateBlocked)
	Msg := "Repository successfully blocked!"
	logger.Success("ActionsController/Block", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONSuccess(Msg, nil)))
}

// Activate describes activate repository action.
func (ctrl ActionsController) Activate(c *gin.Context) {
	var activateActionRequest model.ActivateActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &activateActionRequest)
	if err != nil {
		logger.Error("ActionsController/Activate", err.Error())
		logger.Warning("ActionsController/Activate", "Repository activate aborted!")
		return
	}

	UpdateRepositoryState(activateActionRequest.RepositoryUUID, model.StateActive)
	Msg := "Repository successfully activated!"
	logger.Success("ActionsController/Activate", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONSuccess(Msg, nil)))
}

// BuildWsJSONError returns json formatted error response.
func BuildWsJSONError(message string) string {
	return fmt.Sprintf(`{"status":"error","message":"%s"}`, message)
}

// BuildWsJSONSuccess returns json formatted success response.
func BuildWsJSONSuccess(message string, data interface{}) string {
	if data != nil {
		return fmt.Sprintf(`{"status":"success","message":"%s","data":%v}`, message, data)
	}
	return fmt.Sprintf(`{"status":"success","message":"%s","data":""}`, message)
}
