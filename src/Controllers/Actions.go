package Controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Shitovdm/git-rsync/src/Components/Cmd"
	"github.com/Shitovdm/git-rsync/src/Components/Configuration"
	"github.com/Shitovdm/git-rsync/src/Components/Helpers"
	"github.com/Shitovdm/git-rsync/src/Components/Logger"
	"github.com/Shitovdm/git-rsync/src/Models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ActionsController struct{}

func (ctrl ActionsController) Pull(c *gin.Context) {

	var pullActionRequest Models.PullActionRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &pullActionRequest)
	if err != nil {
		UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, STATUS_PULLFAILED)
		Logger.Error("ActionsController/Pull", err.Error())
		Logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, STATUS_PENDINGPULL)
	Logger.Trace("ActionsController/Pull", fmt.Sprintf("Start pulling repository with UUID %s...", pullActionRequest.RepositoryUuid))

	repositoryConfig := Configuration.GetRepositoryByUuid(pullActionRequest.RepositoryUuid)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, STATUS_PULLFAILED)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", pullActionRequest.RepositoryUuid)
		Logger.Error("ActionsController/Pull", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		Logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	platformConfig := Configuration.GetPlatformByUuid(repositoryConfig.SourcePlatformUuid)
	if platformConfig == nil {
		UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, STATUS_PULLFAILED)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.SourcePlatformUuid)
		Logger.Error("ActionsController/Pull", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		Logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	repositoryName := Configuration.GetRepositorySourceRepositoryName(repositoryConfig)

	isNewRepository := false
	if !Helpers.IsDirExists(Configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name))) ||
		!Helpers.IsDirExists(Configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s\source\%s`, repositoryConfig.Name, repositoryName))) {
		Logger.Trace("ActionsController/Pull", fmt.Sprintf("Repository %s has not been initialized earlier! Initialization...", repositoryConfig.Name))
		err = Helpers.CreateNewDir(Configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s\source`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, STATUS_PULLFAILED)
			ErrorMsg := fmt.Sprintf(`Error while creating new folder .\projects\%s\source`, repositoryConfig.Name)
			Logger.Error("ActionsController/Pull", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
			Logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
			return
		}
		Logger.Trace("ActionsController/Pull", fmt.Sprintf("Root folders for repository %s succesfully created!", repositoryConfig.Name))
		isNewRepository = true
	}

	repositoryFullURL := platformConfig.Address + repositoryConfig.SourcePlatformPath
	repositoryFullPath := Configuration.BuildPlatformPath(fmt.Sprintf(`projects\%s\source`, repositoryConfig.Name))
	//
	finishFetching := make(chan bool)
	if isNewRepository {
		//	Clone action.
		Logger.Trace("ActionsController/Pull", fmt.Sprintf("Cloning repository from %s...", repositoryFullURL))
		go func() {
			cloneResult := Cmd.Clone(repositoryFullPath, repositoryFullURL)
			if cloneResult {
				Logger.Trace("ActionsController/Pull", fmt.Sprintf("Repository %s cloned successfully!", repositoryFullURL))
				finishFetching <- true
			} else {
				Logger.Error("ActionsController/Pull", fmt.Sprintf("Error occurred while cloning repository %s!", repositoryFullURL))
				finishFetching <- false
			}
		}()
	} else {
		//	Pull action.
		Logger.Trace("ActionsController/Pull", fmt.Sprintf("Fetching new from %s...", repositoryFullURL))
		go func() {
			pullResult := Cmd.Pull(repositoryFullPath + `\` + repositoryName)
			if pullResult {
				Logger.Trace("ActionsController/Pull", fmt.Sprintf("Repository %s pulled successfully!", repositoryFullURL))
				finishFetching <- true
			} else {
				Logger.Error("ActionsController/Pull", fmt.Sprintf("Error occurred while pulling repository %s!", repositoryFullURL))
				finishFetching <- false
			}
		}()
	}

	fetchRes := <-finishFetching
	if !fetchRes {
		UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, STATUS_PULLFAILED)
		Msg := "Error occurred while fetching source repository!"
		Logger.Error("ActionsController/Pull", Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(Msg)))
		Logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
		return
	}

	UpdateRepositoryStatus(pullActionRequest.RepositoryUuid, STATUS_PULLED)
	Msg := "Source repository fetched successfully!"
	Logger.Success("ActionsController/Pull", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	return
}

func (ctrl ActionsController) Push(c *gin.Context) {

	var pushActionRequest Models.PushActionRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &pushActionRequest)
	if err != nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_PUSHFAILED)
		Logger.Error("ActionsController/Push", err.Error())
		Logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_PENDINGPUSH)
	Logger.Trace("ActionsController/Push", fmt.Sprintf("Start pushing repository with UUID %s...", pushActionRequest.RepositoryUuid))

	repositoryConfig := Configuration.GetRepositoryByUuid(pushActionRequest.RepositoryUuid)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_PUSHFAILED)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", pushActionRequest.RepositoryUuid)
		Logger.Error("ActionsController/Push", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		Logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	platformConfig := Configuration.GetPlatformByUuid(repositoryConfig.DestinationPlatformUuid)
	if platformConfig == nil {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_PUSHFAILED)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.DestinationPlatformUuid)
		Logger.Error("ActionsController/Push", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		Logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	//	Step 1. Cloning/pulling destination repository.
	destinationRepositoryName := Configuration.GetRepositoryDestinationRepositoryName(repositoryConfig)
	sourceRepositoryName := Configuration.GetRepositorySourceRepositoryName(repositoryConfig)

	isNewRepository := false
	if !Helpers.IsDirExists(Configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name))) ||
		!Helpers.IsDirExists(Configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s\destination\%s`, repositoryConfig.Name, destinationRepositoryName))) {
		Logger.Trace("ActionsController/Push", fmt.Sprintf("Repository %s has not been initialized earlier! Initializing...", repositoryConfig.Name))
		err = Helpers.CreateNewDir(Configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s\destination`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_PUSHFAILED)
			ErrorMsg := fmt.Sprintf("Error while creating new folder ./projects/%s/destination", repositoryConfig.Name)
			Logger.Error("ActionsController/Push", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
			Logger.Warning("ActionsController/Push", "Remote repository update aborted!")
			return
		}
		Logger.Trace("ActionsController/Push", fmt.Sprintf("Root folder for repository %s successfully created!", repositoryConfig.Name))
		isNewRepository = true
	}

	repositoryFullURL := platformConfig.Address + repositoryConfig.DestinationPlatformPath
	repositoryFullPath := Configuration.BuildPlatformPath(fmt.Sprintf(`projects\%s`, repositoryConfig.Name))
	finishFetching := make(chan bool)
	if isNewRepository {
		//	Clone action.
		Logger.Trace("ActionsController/Push", fmt.Sprintf("Cloning repository from %s...", repositoryFullURL))
		go func() {
			cloneResult := Cmd.Clone(repositoryFullPath+`\destination\`, repositoryFullURL)
			if cloneResult {
				Logger.Trace("ActionsController/Push", fmt.Sprintf("Repository %s cloned successfully!", repositoryFullURL))
				finishFetching <- true
			} else {
				Logger.Error("ActionsController/Push", fmt.Sprintf("Error occurred while cloning repository %s!", repositoryFullURL))
				finishFetching <- false
			}
		}()
	} else {
		//	Pull action.
		Logger.Trace("ActionsController/Push", fmt.Sprintf("Pulling new from %s...", repositoryFullURL))
		go func() {
			pullResult := Cmd.Pull(repositoryFullPath + `\destination\` + destinationRepositoryName)
			if pullResult {
				Logger.Trace("ActionsController/Push", fmt.Sprintf("Repository %s pulled successfully!", repositoryFullURL))
				finishFetching <- true
			} else {
				Logger.Error("ActionsController/Push", fmt.Sprintf("Error occurred while pulling repository %s!", repositoryFullURL))
				finishFetching <- false
			}
		}()
	}

	fetchRes := <-finishFetching
	if !fetchRes {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_PUSHFAILED)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError("Error occurred while fetching destination repository!")))
		Logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}

	//	Step 2. Copying all files from source to destination repositories dir`s.
	Logger.Trace("ActionsController/Push", "Start copying all repository files...")
	if !Cmd.CopyRepository(repositoryFullPath, destinationRepositoryName, sourceRepositoryName) {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_PUSHFAILED)
		Logger.Error("ActionsController/Push", "Error occurred while copying repository files!")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError("Error occurred while copying repository files!")))
		Logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}
	Logger.Trace("ActionsController/Push", "Repository files successfully copied!")

	destinationRepositoryPath := repositoryFullPath + `\destination\` + destinationRepositoryName

	//	Step 3. Checking needed pushing destination repository.
	commits, err := Cmd.Log(destinationRepositoryPath, "origin/master..HEAD", -1)
	if err == nil {
		if len(commits) == 0 {
			UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_SYNCHRONIZED)
			Logger.Debug("ActionsController/Push", "Destination repository does not need to be updated, all changes are pushed earlier!")
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess("Destination repository does not need to be updated!", nil)))
			return
		}
	}
	Logger.Trace("ActionsController/Push", "Remote destination repository needs updating, have unpushed changes!")

	//	Step 4. Rewriting commits author (if needed).
	appConfig := Configuration.GetAppConfig()
	if appConfig.CommitsOverriding.State {
		Logger.Trace("ActionsController/Pull", "Overriding source repository commits author...")
		if !Cmd.OverrideAuthor(destinationRepositoryPath, appConfig.CommitsOverriding.CommittersRules) {
			UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_PUSHFAILED)
			Msg := "Error occurred while overriding destination repository commits author!"
			Logger.Error("ActionsController/Pull", Msg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(Msg)))
			Logger.Warning("ActionsController/Pull", "Fetching source repository aborted!")
			return
		}
		Logger.Trace("ActionsController/Pull", "All commits in source repository successfully overridden!")
	}else{
		Logger.Trace("ActionsController/Pull", "Overriding source repository commits not needed!")
	}

	//	Step 5. Pushing destination repository to remote.
	/*Logger.Trace("ActionsController/Push", "Pushing destination repository...")
	if !Cmd.Push(destinationRepositoryPath) {
		UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_PUSHFAILED)
		Logger.Error("ActionsController/Push", "Error occurred while pushing destination repository!")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError("Error occurred while pushing destination repository!")))
		Logger.Warning("ActionsController/Push", "Remote repository update aborted!")
		return
	}*/

	UpdateRepositoryStatus(pushActionRequest.RepositoryUuid, STATUS_SYNCHRONIZED)
	Msg := "Destination repository successfully pushed!"
	Logger.Success("ActionsController/Push", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	return
}

func (ctrl ActionsController) Clear(c *gin.Context) {

	var cleanActionRequest Models.CleanActionRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &cleanActionRequest)
	if err != nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, STATUS_CLEANFAILED)
		Logger.Error("ActionsController/Clear", err.Error())
		Logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
		return
	}

	UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, STATUS_PENDINGCLEAN)
	Logger.Trace("ActionsController/Clear", fmt.Sprintf("Start cleaning repository with UUID %s...", cleanActionRequest.RepositoryUuid))

	repositoryConfig := Configuration.GetRepositoryByUuid(cleanActionRequest.RepositoryUuid)
	if repositoryConfig == nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, STATUS_CLEANFAILED)
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", cleanActionRequest.RepositoryUuid)
		Logger.Error("ActionsController/Clear", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		Logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
		return
	}

	platformConfig := Configuration.GetPlatformByUuid(repositoryConfig.SourcePlatformUuid)
	if platformConfig == nil {
		UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, STATUS_CLEANFAILED)
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.SourcePlatformUuid)
		Logger.Error("ActionsController/Clear", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		Logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
		return
	}

	if Helpers.IsDirExists(Configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name))) {
		err = Helpers.RemoveDir(Configuration.BuildPlatformPath(fmt.Sprintf(`\projects\%s`, repositoryConfig.Name)))
		if err != nil {
			UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, STATUS_CLEANFAILED)
			ErrorMsg := "Error deleting project folder! Check if the folder is used by other programs and try again."
			Logger.Error("ActionsController/Clear", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
			Logger.Warning("ActionsController/Clear", "Repository runtime data cleaning aborted!")
			return
		}
	}

	UpdateRepositoryStatus(cleanActionRequest.RepositoryUuid, STATUS_CLEANED)
	Msg := "Repository runtime data successfully cleaned!"
	Logger.Success("ActionsController/Clear", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	return
}

func (ctrl ActionsController) Info(c *gin.Context) {
	var infoActionRequest Models.InfoActionRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &infoActionRequest)
	if err != nil {
		Logger.Error("ActionsController/Info", err.Error())
		Logger.Warning("ActionsController/Info", "Getting repository info aborted!")
		return
	}

	repositoryConfig := Configuration.GetRepositoryByUuid(infoActionRequest.RepositoryUuid)
	if repositoryConfig == nil {
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", infoActionRequest.RepositoryUuid)
		Logger.Error("ActionsController/Info", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		Logger.Warning("ActionsController/Info", "Getting repository info aborted!")
		return
	}

	destinationRepositoryName := Configuration.GetRepositoryDestinationRepositoryName(repositoryConfig)
	repositoryFullPath := Configuration.BuildPlatformPath(fmt.Sprintf(`projects\%s`, repositoryConfig.Name))
	destinationRepositoryPath := repositoryFullPath + `\destination\` + destinationRepositoryName
	commitsLimit := Configuration.GetAppConfigField("Common", "RecentCommitsShown")
	commits, err := Cmd.Log(destinationRepositoryPath, "", int(commitsLimit.Int()))
	if err != nil {
		ErrorMsg := fmt.Sprintf("Unable to select commits for repository %s!", destinationRepositoryName)
		Logger.Error("ActionsController/Info", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(ErrorMsg)))
		Logger.Warning("ActionsController/Info", "Getting repository info aborted!")
		return
	}

	commitsJson, _ := json.Marshal(commits)
	Msg := "Repository commits list selected successfully!"
	Logger.Success("ActionsController/Info", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, string(commitsJson))))
	return
}

func (ctrl ActionsController) Block(c *gin.Context) {
	var blockActionRequest Models.BlockActionRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &blockActionRequest)
	if err != nil {
		Logger.Error("ActionsController/Block", err.Error())
		Logger.Warning("ActionsController/Block", "Repository blocking aborted!")
		return
	}

	UpdateRepositoryState(blockActionRequest.RepositoryUuid, Models.STATE_BLOCKED)
	Msg := "Repository successfully blocked!"
	Logger.Success("ActionsController/Block", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	return
}

func (ctrl ActionsController) Activate(c *gin.Context) {
	var activateActionRequest Models.ActivateActionRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &activateActionRequest)
	if err != nil {
		Logger.Error("ActionsController/Activate", err.Error())
		Logger.Warning("ActionsController/Activate", "Repository activate aborted!")
		return
	}

	UpdateRepositoryState(activateActionRequest.RepositoryUuid, Models.STATE_ACTIVE)
	Msg := "Repository successfully activated!"
	Logger.Success("ActionsController/Activate", Msg)
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
