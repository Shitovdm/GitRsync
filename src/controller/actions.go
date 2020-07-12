package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/actions"
	"github.com/Shitovdm/GitRsync/src/component/cmd"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/Shitovdm/GitRsync/src/model"
	"github.com/Shitovdm/GitRsync/src/model/repository"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ActionsController struct describes main actions section controller.
type ActionsController struct{}

// Pull describes pull repository action.
func (ctrl ActionsController) Pull(c *gin.Context) {

	var pullActionRequest model.PullActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &pullActionRequest)
	repo := repository.Get(pullActionRequest.RepositoryUUID)
	if err != nil {
		repo.SetStatus(repository.StatusPullFailed)
		_ = repo.Update()
		LogErrorWithoutResponse("Pull", "Fetching source repository aborted!", err)
		return
	}


	err = actions.PullCommits(pullActionRequest.RepositoryUUID)
	if err != nil {
		repo.SetStatus(repository.StatusPullFailed)
		_ = repo.Update()
		LogErrorWithResponse("Pull", err.Error(), "Pulling source repository aborted!", conn)
		return
	}

	repo.SetStatus(repository.StatusPulled)
	_ = repo.Update()
	LogSuccessWithResponse("Pull", "Source repository fetched successfully!", nil, conn)
}

// Push describes push repository action.
func (ctrl ActionsController) Push(c *gin.Context) {

	var pushActionRequest model.PushActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &pushActionRequest)
	if err != nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusPushFailed)
		logger.Error("ActionsController/Push", err.Error())
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusPendingPush)
	logger.Trace("ActionsController/Push", fmt.Sprintf("Start pushing repository with UUID %s...", pushActionRequest.RepositoryUUID))

	repositoryConfig := conf.GetRepositoryByUUID(pushActionRequest.RepositoryUUID)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusPushFailed)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", pushActionRequest.RepositoryUUID)
		logger.Error("ActionsController/Push", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(ErrorMsg)))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	platformConfig := conf.GetPlatformByUUID(repositoryConfig.DestinationPlatformUUID)
	if platformConfig == nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusPushFailed)
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
			UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusPushFailed)
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
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusPushFailed)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError("Error occurred while fetching destination repository!")))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	//	Step 2. Copying all files from source to destination repositories dir`s.
	logger.Trace("ActionsController/Push", "Start copying all repository files...")
	if !cmd.CopyRepository(repositoryFullPath, destinationRepositoryName, sourceRepositoryName) {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusPushFailed)
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
			UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusSynchronized)
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
			UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusPushFailed)
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
		UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusPushFailed)
		logger.Error("ActionsController/Push", "Error occurred while pushing destination repository!")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError("Error occurred while pushing destination repository!")))
		logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	UpdateRepositoryStatus(pushActionRequest.RepositoryUUID, repository.StatusSynchronized)
	Msg := "Destination repository successfully pushed!"
	logger.Success("ActionsController/Push", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONSuccess(Msg, nil)))
}

// Clear describes clear repository data action.
func (ctrl ActionsController) Clear(c *gin.Context) {

	var cleanActionRequest model.CleanActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &cleanActionRequest)
	if err != nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, repository.StatusCleanFailed)
		LogErrorWithoutResponse("Clear", "Repository runtime data cleaning aborted!", err)
		return
	}

	UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, repository.StatusPendingClean)
	logger.Trace("ActionsController/Clear", fmt.Sprintf("Start cleaning repository with UUID %s...", cleanActionRequest.RepositoryUUID))

	repositoryConfig := conf.GetRepositoryByUUID(cleanActionRequest.RepositoryUUID)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, repository.StatusCleanFailed)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", cleanActionRequest.RepositoryUUID)
		LogErrorWithResponse("Clear", ErrorMsg, "Repository runtime data cleaning aborted!", conn)
		return
	}

	platformConfig := conf.GetPlatformByUUID(repositoryConfig.SourcePlatformUUID)
	if platformConfig == nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, repository.StatusCleanFailed)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.SourcePlatformUUID)
		LogErrorWithResponse("Clear", ErrorMsg, "Repository runtime data cleaning aborted!", conn)
		return
	}

	if helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name))) {
		err = helper.RemoveDir(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, repository.StatusCleanFailed)
			ErrorMsg := "Error deleting project folder! Check if the folder is used by other programs and try again."
			LogErrorWithResponse("Clear", ErrorMsg, "Repository runtime data cleaning aborted!", conn)
			return
		}
	}

	UpdateRepositoryStatus(cleanActionRequest.RepositoryUUID, repository.StatusCleaned)
	LogSuccessWithResponse("Clear", "Repository runtime data successfully cleaned!", nil, conn)
}

// Info describes repository info action.
func (ctrl ActionsController) Info(c *gin.Context) {
	var infoActionRequest model.InfoActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &infoActionRequest)
	if err != nil {
		LogErrorWithoutResponse("Info", "Getting repository info aborted!", err)
		return
	}

	repositoryConfig := conf.GetRepositoryByUUID(infoActionRequest.RepositoryUUID)
	if repositoryConfig == nil {
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", infoActionRequest.RepositoryUUID)
		LogErrorWithResponse("Info", ErrorMsg, "Getting repository info aborted!", conn)
		return
	}

	destinationRepositoryName := conf.GetRepositoryDestinationRepositoryName(repositoryConfig)
	repositoryFullPath := conf.BuildPlatformPath(fmt.Sprintf(`projects\%s`, repositoryConfig.Name))
	destinationRepositoryPath := repositoryFullPath + `\destination\` + destinationRepositoryName
	commitsLimit := conf.GetAppConfigField("Common", "RecentCommitsShown")
	commits, err := cmd.Log(destinationRepositoryPath, "", int(commitsLimit.Int()))
	if err != nil {
		UpdateRepositoryStatus(infoActionRequest.RepositoryUUID, repository.StatusFailed)
		LogErrorWithResponse("Info", err.Error(), "Getting repository commits aborted!", conn)
		return
	}

	commitsJSON, _ := json.Marshal(commits)
	LogSuccessWithResponse("Info", "Repository commits list selected successfully!", commitsJSON, conn)
}

// Block describes block repository action.
func (ctrl ActionsController) Block(c *gin.Context) {
	var blockActionRequest model.BlockActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &blockActionRequest)
	if err != nil {
		LogErrorWithoutResponse("Block", "Repository blocking aborted!", err)
		return
	}

	UpdateRepositoryState(blockActionRequest.RepositoryUUID, repository.StateBlocked)
	LogSuccessWithResponse("Block", "Repository successfully blocked!", nil, conn)
}

// Activate describes activate repository action.
func (ctrl ActionsController) Activate(c *gin.Context) {
	var activateActionRequest model.ActivateActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &activateActionRequest)
	if err != nil {
		LogErrorWithoutResponse("Activate", "Repository activate aborted!", err)
		return
	}

	UpdateRepositoryState(activateActionRequest.RepositoryUUID, repository.StateActive)
	LogSuccessWithResponse("Activate", "Repository successfully activated!", nil, conn)
}

// OpenDir opens fs dir in explorer.
func (ctrl ActionsController) OpenDir(c *gin.Context) {
	var openDirActionRequest model.OpenDirActionRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &openDirActionRequest)
	if err != nil {
		LogErrorWithoutResponse("OpenDir", "Opening file system folder aborted!", err)
		return
	}

	logger.Success("ActionsController/OpenDir", "Opening file system folder "+openDirActionRequest.Path+" in explorer...")
	helper.ExploreDir(conf.BuildPlatformPath(openDirActionRequest.Path))
}

// SyncTags syncs repositories tags.
func (ctrl ActionsController) SyncTags(c *gin.Context) {

	var syncTagsActionRequest model.SyncTagsActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &syncTagsActionRequest)
	if err != nil {
		LogErrorWithoutResponse("SyncTags", "Syncing repositories tags aborted!", err)
		return
	}

	err = actions.SyncTags(syncTagsActionRequest.RepositoryUUID)
	if err != nil {
		LogErrorWithResponse("SyncTags", err.Error(), "Getting repository info aborted!", conn)
		return
	}

	LogSuccessWithResponse("SyncTags", "Repository tags successfully synchronized!", nil, conn)
}

// LogErrorWithResponse logs error with response.
func LogErrorWithResponse(action, errMessage, operationErrMessage string, conn *websocket.Conn) {
	logger.Error("ActionsController/"+action, errMessage)
	errReply := fmt.Sprintf(`{"status":"error","message":"%s"}`, errMessage)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(errReply))
	logger.Warning("ActionsController/"+action, operationErrMessage)
}

// LogSuccessWithResponse logs success with response.
func LogSuccessWithResponse(action, message string, data interface{}, conn *websocket.Conn) {
	logger.Success("ActionsController/"+action, message)
	reply := fmt.Sprintf(`{"status":"success","message":"%s","data":""}`, message)
	if data != nil {
		reply = fmt.Sprintf(`{"status":"success","message":"%s","data":%v}`, message, data)
	}
	_ = conn.WriteMessage(websocket.TextMessage, []byte(reply))
}

// LogErrorWithoutResponse logs error without response.
func LogErrorWithoutResponse(action, operationErrMessage string, err error) {
	logger.Error("ActionsController/"+action, err.Error())
	logger.Warning("ActionsController/"+action, operationErrMessage)
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
