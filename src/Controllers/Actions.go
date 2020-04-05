package Controllers

import (
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Cmd"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Logger"
	"github.com/Shitovdm/git-repo-exporter/src/Models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

type ActionsController struct{}

func (ctrl ActionsController) Sync(c *gin.Context) {

}

func (ctrl ActionsController) Pull(c *gin.Context) {

	var pullActionRequest Models.PullActionRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &pullActionRequest)
	if err != nil {
		Logger.Error("ActionsController/Pull", err.Error())
		return
	}
	Logger.Trace("ActionsController/Pull", fmt.Sprintf("Start pulling repository with UUID %s...", pullActionRequest.RepositoryUuid))

	repositoryConfig := Configuration.GetRepositoryByUuid(pullActionRequest.RepositoryUuid)
	if repositoryConfig == nil {
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", pullActionRequest.RepositoryUuid)
		Logger.Error("ActionsController/Pull", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(ErrorMsg))
		return
	}

	platformConfig := Configuration.GetPlatformByUuid(repositoryConfig.SourcePlatformUuid)
	if platformConfig == nil {
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.SourcePlatformUuid)
		Logger.Error("ActionsController/Pull", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(ErrorMsg))
		return
	}

	spp := strings.TrimRight(repositoryConfig.SourcePlatformPath, ".git")
	repositoryName := strings.Split(spp, "/")[len(strings.Split(spp, "/"))-1]

	isNewRepository := false
	if !Helpers.IsDirExists(Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s", repositoryConfig.Name))) ||
		!Helpers.IsDirExists(Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s/%s", repositoryConfig.Name, repositoryName))) {
		Logger.Info("ActionsController/Pull", fmt.Sprintf("Repository %s has not been initialized earlier! Initialization...", repositoryConfig.Name))
		err = Helpers.CreateNewDir(Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s", repositoryConfig.Name)))
		if err != nil {
			ErrorMsg := fmt.Sprintf("Error while creating new folder ./projects/%s", repositoryConfig.Name)
			Logger.Error("ActionsController/Pull", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(ErrorMsg))
			return
		}
		Logger.Info("ActionsController/Pull", fmt.Sprintf("Root folders for repository %s succesfully created!", repositoryConfig.Name))
		isNewRepository = true
	}

	repositoryFullURL := platformConfig.Address + repositoryConfig.SourcePlatformPath
	repositoryFullPath := Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s", repositoryConfig.Name))
	if isNewRepository {
		//	Clone action.
		Logger.Trace("ActionsController/Pull", fmt.Sprintf("Cloning repository from %s...", repositoryFullURL))
		go func() {
			cloneResult := Cmd.Clone(repositoryFullPath, repositoryFullURL)
			if cloneResult {
				Logger.Info("ActionsController/Pull", fmt.Sprintf("Repository %s cloned successfully!", repositoryFullURL))
				_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess("Repository cloned successfully!")))
			} else {
				Logger.Error("ActionsController/Pull", fmt.Sprintf("Error occurred while cloning repository %s!", repositoryFullURL))
				_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError("Error occurred while cloning the repository!")))
			}
			_ = conn.Close()
		}()
	} else {
		//	Pull action.
		Logger.Trace("ActionsController/Pull", fmt.Sprintf("Fetching new from %s...", repositoryFullURL))
		go func() {
			pullResult := Cmd.Pull(repositoryFullPath + "/" + repositoryName)
			if pullResult {
				Logger.Info("ActionsController/Pull", fmt.Sprintf("Repository %s pulled successfully!", repositoryFullURL))
				_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess("Repository pulled successfully!")))
			} else {
				Logger.Error("ActionsController/Pull", fmt.Sprintf("Error occurred while pulling repository %s!", repositoryFullURL))
				_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError("Error occurred while pulling repository!")))
			}
			_ = conn.Close()
		}()
	}
}

func (ctrl ActionsController) Push(c *gin.Context) {

	var pushActionRequest Models.PushActionRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &pushActionRequest)
	if err != nil {
		Logger.Error("ActionsController/Push", err.Error())
		return
	}
	Logger.Trace("ActionsController/Push", fmt.Sprintf("Start pushing repository with UUID %s...", pushActionRequest.RepositoryUuid))

	repositoryConfig := Configuration.GetRepositoryByUuid(pushActionRequest.RepositoryUuid)
	if repositoryConfig == nil {
		ErrorMsg := fmt.Sprintf("Repository with transferred UUID %s not found!", pushActionRequest.RepositoryUuid)
		Logger.Error("ActionsController/Push", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(ErrorMsg))
		return
	}

	platformConfig := Configuration.GetPlatformByUuid(repositoryConfig.DestinationPlatformUuid)
	if platformConfig == nil {
		ErrorMsg := fmt.Sprintf("Platform with UUID %s not found!", repositoryConfig.DestinationPlatformUuid)
		Logger.Error("ActionsController/Push", ErrorMsg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(ErrorMsg))
		return
	}

	//	Step 1. Cloning/pulling destination repository.
	dpp := strings.TrimRight(repositoryConfig.DestinationPlatformPath, ".git")
	destinationRepositoryName := strings.Split(dpp, "/")[len(strings.Split(dpp, "/"))-1]

	spp := strings.TrimRight(repositoryConfig.SourcePlatformPath, ".git")
	sourceRepositoryName := strings.Split(spp, "/")[len(strings.Split(spp, "/"))-1]

	isNewRepository := false
	if !Helpers.IsDirExists(Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s", repositoryConfig.Name))) ||
		!Helpers.IsDirExists(Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s/%s", repositoryConfig.Name, destinationRepositoryName))) {
		Logger.Info("ActionsController/Push", fmt.Sprintf("Repository %s has not been initialized earlier! Creating root folder for repository...", repositoryConfig.Name))
		err = Helpers.CreateNewDir(Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s", repositoryConfig.Name)))
		if err != nil {
			ErrorMsg := fmt.Sprintf("Error while creating new folder ./projects/%s", repositoryConfig.Name)
			Logger.Error("ActionsController/Push", ErrorMsg)
			_ = conn.WriteMessage(websocket.TextMessage, []byte(ErrorMsg))
			return
		}
		Logger.Trace("ActionsController/Push", "Done!")
		isNewRepository = true
	}

	repositoryFullURL := platformConfig.Address + repositoryConfig.DestinationPlatformPath
	repositoryFullPath := Configuration.BuildPlatformPath(fmt.Sprintf("/projects/%s", repositoryConfig.Name))
	finishFetching := make(chan bool)
	if isNewRepository {
		//	Clone action.
		Logger.Trace("ActionsController/Push", fmt.Sprintf("Cloning repository from %s...", repositoryFullURL))
		go func() {
			cloneResult := Cmd.Clone(repositoryFullPath, repositoryFullURL)
			if cloneResult {
				Logger.Info("ActionsController/Push", fmt.Sprintf("Repository %s cloned successfully!", repositoryFullURL))
				finishFetching<- true
			} else {
				Logger.Error("ActionsController/Push", fmt.Sprintf("Error occurred while cloning repository %s!", repositoryFullURL))
				finishFetching<- false
			}
		}()
	} else {
		//	Pull action.
		Logger.Trace("ActionsController/Push", fmt.Sprintf("Fetching new from %s...", repositoryFullURL))
		go func() {
			pullResult := Cmd.Pull(repositoryFullPath + "/" + destinationRepositoryName)
			if pullResult {
				Logger.Info("ActionsController/Push", fmt.Sprintf("Repository %s pulled successfully!", repositoryFullURL))
				finishFetching<- true
			} else {
				Logger.Error("ActionsController/Push", fmt.Sprintf("Error occurred while pulling repository %s!", repositoryFullURL))
				finishFetching<- false
			}
		}()
	}

	fetchRes := <-finishFetching
	if !fetchRes {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError("Error occurred while fetching destination repository!")))
		_ = conn.Close()
		return
	}

	//	Step 2. Copying all files from source to destination repositories dir`s.
	Logger.Trace("ActionsController/Push", "Start copying all repository files...")
	if !Cmd.CopyRepository(repositoryFullPath, destinationRepositoryName, sourceRepositoryName) {
		Logger.Error("ActionsController/Push", "Error occurred while copying repository files!")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess("Error occurred while copying repository files!")))
		_ = conn.Close()
		return
	}
	Logger.Trace("ActionsController/Push", "Done!")

	//	Step 3.1. Checking needed pushing destination repository.
	if Cmd.Status(repositoryFullPath + "/" + destinationRepositoryName){
		Logger.Info("ActionsController/Push", "Destination repository does not need to be updated, all changes are pushed earlier!")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess("Destination repository does not need to be updated!")))
		_ = conn.Close()
		return
	}
	Logger.Trace("ActionsController/Push", "Remote destination repository needs updating, have unpushed changes!")

	//	Step 3.2. Rewriting commits author (if needed).


	//	Step 4. Pushing destination repository to remote.
	Logger.Trace("ActionsController/Push", "Pushing destination repository...")
	if !Cmd.Push(repositoryFullPath + "/" + destinationRepositoryName){
		Logger.Error("ActionsController/Push", "Error occurred while pushing destination repository!")
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess("Error occurred while pushing destination repository!")))
		_ = conn.Close()
		return
	}

	Logger.Trace("ActionsController/Push", "Done!")
	Msg := "Destination repository successfully pushed!"
	Logger.Info("ActionsController/Push", Msg)
	_ = conn.WriteMessage(websocket.TextMessage, []byte(Msg))
	return
}

func (ctrl ActionsController) Block(c *gin.Context) {
	var blockActionRequest Models.BlockActionRequest
	err := c.BindJSON(&blockActionRequest)
	if err != nil {
		Logger.Error("ActionsController/Block", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	RespondWithSuccess(c, map[string]interface{}{})
}

func (ctrl ActionsController) Active(c *gin.Context) {
	var activeActionRequest Models.ActiveActionRequest
	err := c.BindJSON(&activeActionRequest)
	if err != nil {
		Logger.Error("ActionsController/Active", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	RespondWithSuccess(c, map[string]interface{}{})
}

func (ctrl ActionsController) Info(c *gin.Context) {

}

func BuildWsJsonError(message string) string {
	return `{"status":"error","message":"` + message + `"}`
}

func BuildWsJsonSuccess(message string) string {
	return `{"status":"success","message":"` + message + `"}`
}

func RespondWithError(c *gin.Context, message string) {
	c.JSON(200, gin.H{
		"status":  "error",
		"message": message,
	})
}

func RespondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"status": "success",
		"data":   data,
	})
}
