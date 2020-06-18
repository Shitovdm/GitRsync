package Controllers

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/Components/Configuration"
	"github.com/Shitovdm/GitRsync/src/Components/Helpers"
	"github.com/Shitovdm/GitRsync/src/Components/Interface"
	"github.com/Shitovdm/GitRsync/src/Components/Logger"
	"github.com/Shitovdm/GitRsync/src/Models"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

type RepositoriesController struct{}

const (
	StatusInitiated    = "initiated"
	StatusPendingPull  = "pending_pull"
	StatusPulled       = "pulled"
	StatusPullFailed   = "pull_failed"
	StatusPendingPush  = "pending_push"
	StatusPushed       = "pushed"
	StatusPushFailed   = "push_failed"
	StatusPendingClean = "pending_clear"
	StatusCleaned      = "cleared"
	StatusCleanFailed  = "clear_failed"
	StatusSynchronized = "synced"
	StatusFailed       = "failed"
)

var (
	timeFormat = "02-01-2006 15:04"
)

func (ctrl RepositoriesController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Repositories"
	templateParams["repositories"], _ = Configuration.GetRepositoriesConfigData()
	templateParams["platforms"], _ = Configuration.GetPlatformsConfigData()
	c.HTML(http.StatusOK, "repositories/index", templateParams)
}

func (ctrl RepositoriesController) Add(c *gin.Context) {

	var addRepositoryRequest Models.AddRepositoryRequest
	_, err := Helpers.WsHandler(c.Writer, c.Request, &addRepositoryRequest)
	if err != nil {
		Logger.Error("RepositoriesController/Add", err.Error())
		return
	}

	newRepositoryUuid, _ := uuid.NewV4()
	repositories := Configuration.GetRepositoriesConfig()
	repositories = append(repositories, Models.RepositoryConfig{
		Uuid:                    newRepositoryUuid.String(),
		Name:                    addRepositoryRequest.Name,
		SourcePlatformUuid:      addRepositoryRequest.SourcePlatformUuid,
		SourcePlatformPath:      addRepositoryRequest.SourcePlatformPath,
		DestinationPlatformUuid: addRepositoryRequest.DestinationPlatformUuid,
		DestinationPlatformPath: addRepositoryRequest.DestinationPlatformPath,
		Status:                  StatusInitiated,
		State:                   "active",
		UpdatedAt:               "",
	})

	err = Configuration.SaveRepositoriesConfig(repositories)
	if err != nil {
		Logger.Error("RepositoriesController/Add", err.Error())
	}

	Logger.Info("RepositoriesController/Add", fmt.Sprintf("New repository with name %s added successfully!", addRepositoryRequest.Name))
	return
}

func (ctrl RepositoriesController) Edit(c *gin.Context) {

	var editRepositoryRequest Models.EditRepositoryRequest
	_, err := Helpers.WsHandler(c.Writer, c.Request, &editRepositoryRequest)
	if err != nil {
		Logger.Error("RepositoriesController/Edit", err.Error())
		return
	}

	oldRepositoriesList := Configuration.GetRepositoriesConfig()
	newRepositoriesList := make([]Models.RepositoryConfig, 0)
	for _, repository := range oldRepositoriesList {
		if repository.Uuid == editRepositoryRequest.Uuid {
			newRepositoriesList = append(newRepositoriesList, Models.RepositoryConfig{
				Uuid:                    repository.Uuid,
				Name:                    editRepositoryRequest.Name,
				SourcePlatformUuid:      editRepositoryRequest.SourcePlatformUuid,
				SourcePlatformPath:      editRepositoryRequest.SourcePlatformPath,
				DestinationPlatformUuid: editRepositoryRequest.DestinationPlatformUuid,
				DestinationPlatformPath: editRepositoryRequest.DestinationPlatformPath,
				Status:                  repository.Status,
				State:                   repository.State,
				UpdatedAt:               repository.UpdatedAt,
			})
			continue
		}
		newRepositoriesList = append(newRepositoriesList, repository)
	}

	err = Configuration.SaveRepositoriesConfig(newRepositoriesList)
	if err != nil {
		Logger.Error("RepositoriesController/Edit", err.Error())
	}

	Logger.Info("RepositoriesController/Edit", fmt.Sprintf("Repository with name %s successfully edited!", editRepositoryRequest.Name))
	return
}

func (ctrl RepositoriesController) Remove(c *gin.Context) {

	var removeRepositoryRequest Models.RemoveRepositoryRequest
	_, err := Helpers.WsHandler(c.Writer, c.Request, &removeRepositoryRequest)
	if err != nil {
		Logger.Error("RepositoriesController/Remove", err.Error())
		return
	}

	removedRepositoryName := ""
	oldRepositoriesList := Configuration.GetRepositoriesConfig()
	newRepositoriesList := make([]Models.RepositoryConfig, 0)
	for _, repository := range oldRepositoriesList {
		if repository.Uuid != removeRepositoryRequest.Uuid {
			newRepositoriesList = append(newRepositoriesList, repository)
		} else {
			removedRepositoryName = repository.Name
		}
	}

	err = Configuration.SaveRepositoriesConfig(newRepositoriesList)
	if err != nil {
		Logger.Error("RepositoriesController/Remove", err.Error())
	}

	Logger.Info("RepositoriesController/Remove", fmt.Sprintf("Repository with name %s successfully removed!", removedRepositoryName))
	return
}

func UpdateRepositoryStatus(uuid string, status string) {

	oldRepositoriesList := Configuration.GetRepositoriesConfig()
	for i, repository := range oldRepositoriesList {
		if repository.Uuid == uuid {
			oldRepositoriesList[i].Status = status
			if status == StatusPulled || status == StatusPushed || status == StatusSynchronized {
				t := time.Now()
				oldRepositoriesList[i].UpdatedAt = t.Format(timeFormat)
			}
		}
	}

	err := Configuration.SaveRepositoriesConfig(oldRepositoriesList)
	if err != nil {
		return
	}

	return
}

func UpdateRepositoryState(uuid string, state string) {

	oldRepositoriesList := Configuration.GetRepositoriesConfig()
	for i, repository := range oldRepositoriesList {
		if repository.Uuid == uuid {
			oldRepositoriesList[i].State = state
			t := time.Now()
			oldRepositoriesList[i].UpdatedAt = t.Format(timeFormat)
		}
	}

	err := Configuration.SaveRepositoriesConfig(oldRepositoriesList)
	if err != nil {
		return
	}

	return
}
