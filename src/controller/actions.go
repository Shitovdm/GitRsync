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
	"github.com/Shitovdm/GitRsync/src/model/platform"
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
		LogErrorWithoutResponse("Pull", "Pulling source repository aborted!", err)
		return
	}

	err = actions.PullCommits(repo, "source")
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
	repo := repository.Get(pushActionRequest.RepositoryUUID)
	if err != nil {
		repo.SetStatus(repository.StatusPushFailed)
		_ = repo.Update()
		LogErrorWithResponse("Push", err.Error(), "Pushing destination repository aborted!", conn)
		return
	}
	repo.SetStatus(repository.StatusPendingPush)
	_ = repo.Update()

	//	Step 1. Cloning/pulling destination repository.
	err = actions.PullCommits(repo, "source")
	if err != nil {
		repo.SetStatus(repository.StatusPullFailed)
		_ = repo.Update()
		LogErrorWithResponse("Pull", err.Error(), "Pulling source repository aborted!", conn)
		return
	}
	logger.Trace("ActionsController/Push", "Repository successfully fetched!")

	//	Step 2. Copying all files from source to destination repositories dir`s.
	repositoryFullPath := conf.BuildPlatformPath(fmt.Sprintf(`projects\%s\destination`, repo.GetName()))
	sourceName := repo.GetSourceRepositoryName()
	destinationName := repo.GetDestinationRepositoryName()
	logger.Trace("ActionsController/Push", "Start copying all repository files...")
	if !cmd.CopyRepository(repositoryFullPath, destinationName, sourceName) {
		repo.SetStatus(repository.StatusPushFailed)
		_ = repo.Update()
		LogErrorWithResponse("Push", "Error occurred while copying repository files!",
			"Pulling source repository aborted!", conn)
		return
	}
	logger.Trace("ActionsController/Push", "Repository files successfully copied!")

	//	Step 3. Checking needed pushing destination repository.

	commits, err := cmd.Log(repo.GetDestinationRepositoryName(), "origin/master..HEAD", -1)
	if err == nil {
		if len(commits) == 0 {
			repo.SetStatus(repository.StatusSynchronized)
			_ = repo.Update()
			LogSuccessWithResponse("Pull", "Destination repository does not need to be updated!", nil, conn)
			return
		}
	}
	logger.Trace("ActionsController/Push", "Remote destination repository needs updating, have unpushed changes!")

	//	Step 4. Rewriting commits author (if needed).
	appConfig := conf.GetAppConfig()
	if appConfig.CommitsOverriding.State {
		logger.Trace("ActionsController/Push", "Overriding source repository commits author...")
		err = actions.OverrideCommits(repo)
		if err != nil {
			repo.SetStatus(repository.StatusPushFailed)
			_ = repo.Update()
			LogErrorWithResponse("Push", err.Error(), "Pushing destination repository aborted!", conn)
			return
		}
		logger.Trace("ActionsController/Push", "All commits in source repository successfully overridden!")
	} else {
		logger.Trace("ActionsController/Push", "Overriding source repository commits not needed!")
	}

	//	Step 5. Pushing destination repository to remote.
	err = actions.PushCommits(repo)
	if err != nil {
		repo.SetStatus(repository.StatusPushFailed)
		_ = repo.Update()
		LogErrorWithResponse("Push", err.Error(), "Pushing destination repository aborted!", conn)
		return
	}

	repo.SetStatus(repository.StatusPushed)
	_ = repo.Update()
	LogSuccessWithResponse("Push", "Destination repository successfully pushed!", nil, conn)
}

// Clear describes clear repository data action.
func (ctrl ActionsController) Clear(c *gin.Context) {

	var cleanActionRequest model.CleanActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &cleanActionRequest)
	repo := repository.Get(cleanActionRequest.RepositoryUUID)
	if err != nil {
		repo.SetStatus(repository.StatusPushFailed)
		_ = repo.Update()
		LogErrorWithResponse("Clear", err.Error(), "Repository runtime data cleaning aborted!", conn)
		return
	}

	repo.SetStatus(repository.StatusPendingClean)
	_ = repo.Update()
	logger.Trace("ActionsController/Clear", fmt.Sprintf("Cleaning repository runtime data with UUID %s...", cleanActionRequest.RepositoryUUID))

	pl := platform.Get(repo.DestinationPlatformUUID)
	if pl == nil {
		repo.SetStatus(repository.StatusPushFailed)
		_ = repo.Update()
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repo.GetSourcePlatformUUID())
		LogErrorWithResponse("Clear", ErrorMsg, "Repository runtime data cleaning aborted!", conn)
		return
	}

	if helper.IsDirExists(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repo.GetName()))) {
		err = helper.RemoveDir(conf.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repo.GetName())))
		if err != nil {
			repo.SetStatus(repository.StatusCleanFailed)
			_ = repo.Update()
			LogErrorWithResponse("Clear", "Error deleting project folder! Check if the folder is used by other programs and try again.",
				"Repository runtime data cleaning aborted!", conn)
			return
		}
	}

	repo.SetStatus(repository.StatusCleaned)
	_ = repo.Update()
	LogSuccessWithResponse("Clear", "Repository runtime data successfully cleaned!", nil, conn)
}

// Info describes repository info action.
func (ctrl ActionsController) Info(c *gin.Context) {
	var infoActionRequest model.InfoActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &infoActionRequest)
	repo := repository.Get(infoActionRequest.RepositoryUUID)
	if err != nil {
		LogErrorWithResponse("Info", err.Error(), "Getting repository info aborted!", conn)
		return
	}

	commitsLimit := conf.GetAppConfigField("Common", "RecentCommitsShown")
	commits, err := actions.GetCommits(repo, int(commitsLimit.Int()))
	if err != nil {
		repo.SetStatus(repository.StatusFailed)
		_ = repo.Update()
		LogErrorWithResponse("Info", "The system cannot find the path specified!", "Getting repository commits aborted!", conn)
		return
	}

	commitsJSON, _ := json.Marshal(commits)
	LogSuccessWithResponse("Info", "Repository commits list selected successfully!", commitsJSON, conn)
}

// Block describes block repository action.
func (ctrl ActionsController) Block(c *gin.Context) {
	var blockActionRequest model.BlockActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &blockActionRequest)
	repo := repository.Get(blockActionRequest.RepositoryUUID)
	if err != nil {
		LogErrorWithoutResponse("Block", "Repository blocking aborted!", err)
		return
	}

	repo.SetState(repository.StateBlocked)
	_ = repo.Update()
	LogSuccessWithResponse("Block", "Repository successfully blocked!", nil, conn)
}

// Activate describes activate repository action.
func (ctrl ActionsController) Activate(c *gin.Context) {
	var activateActionRequest model.ActivateActionRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &activateActionRequest)
	repo := repository.Get(activateActionRequest.RepositoryUUID)
	if err != nil {
		LogErrorWithoutResponse("Activate", "Repository activate aborted!", err)
		return
	}

	repo.SetState(repository.StateActive)
	_ = repo.Update()
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
	repo := repository.Get(syncTagsActionRequest.RepositoryUUID)
	if err != nil {
		LogErrorWithoutResponse("SyncTags", "Syncing repositories tags aborted!", err)
		return
	}

	err = actions.SyncTags(repo)
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
