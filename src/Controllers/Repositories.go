package controllers

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/components/configuration"
	"github.com/Shitovdm/GitRsync/src/components/helpers"
	"github.com/Shitovdm/GitRsync/src/components/interface"
	"github.com/Shitovdm/GitRsync/src/components/logger"
	"github.com/Shitovdm/GitRsync/src/models"
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
	templateParams["repositories"], _ = configuration.GetRepositoriesConfigData()
	templateParams["platforms"], _ = configuration.GetPlatformsConfigData()
	c.HTML(http.StatusOK, "repositories/index", templateParams)
}

func (ctrl RepositoriesController) Add(c *gin.Context) {

	var addRepositoryRequest models.AddRepositoryRequest
	_, err := helpers.WsHandler(c.Writer, c.Request, &addRepositoryRequest)
	if err != nil {
		logger.Error("RepositoriesController/Add", err.Error())
		return
	}

	newRepositoryUuid, _ := uuid.NewV4()
	repositories := configuration.GetRepositoriesConfig()
	repositories = append(repositories, models.RepositoryConfig{
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

	err = configuration.SaveRepositoriesConfig(repositories)
	if err != nil {
		logger.Error("RepositoriesController/Add", err.Error())
	}

	logger.Info("RepositoriesController/Add", fmt.Sprintf("New repository with name %s added successfully!", addRepositoryRequest.Name))
	return
}

func (ctrl RepositoriesController) Edit(c *gin.Context) {

	var editRepositoryRequest models.EditRepositoryRequest
	_, err := helpers.WsHandler(c.Writer, c.Request, &editRepositoryRequest)
	if err != nil {
		logger.Error("RepositoriesController/Edit", err.Error())
		return
	}

	oldRepositoriesList := configuration.GetRepositoriesConfig()
	newRepositoriesList := make([]models.RepositoryConfig, 0)
	for _, repository := range oldRepositoriesList {
		if repository.Uuid == editRepositoryRequest.Uuid {
			newRepositoriesList = append(newRepositoriesList, models.RepositoryConfig{
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

	err = configuration.SaveRepositoriesConfig(newRepositoriesList)
	if err != nil {
		logger.Error("RepositoriesController/Edit", err.Error())
	}

	logger.Info("RepositoriesController/Edit", fmt.Sprintf("Repository with name %s successfully edited!", editRepositoryRequest.Name))
	return
}

func (ctrl RepositoriesController) Remove(c *gin.Context) {

	var removeRepositoryRequest models.RemoveRepositoryRequest
	_, err := helpers.WsHandler(c.Writer, c.Request, &removeRepositoryRequest)
	if err != nil {
		logger.Error("RepositoriesController/Remove", err.Error())
		return
	}

	removedRepositoryName := ""
	oldRepositoriesList := configuration.GetRepositoriesConfig()
	newRepositoriesList := make([]models.RepositoryConfig, 0)
	for _, repository := range oldRepositoriesList {
		if repository.Uuid != removeRepositoryRequest.Uuid {
			newRepositoriesList = append(newRepositoriesList, repository)
		} else {
			removedRepositoryName = repository.Name
		}
	}

	err = configuration.SaveRepositoriesConfig(newRepositoriesList)
	if err != nil {
		logger.Error("RepositoriesController/Remove", err.Error())
	}

	logger.Info("RepositoriesController/Remove", fmt.Sprintf("Repository with name %s successfully removed!", removedRepositoryName))
	return
}

func UpdateRepositoryStatus(uuid string, status string) {

	oldRepositoriesList := configuration.GetRepositoriesConfig()
	for i, repository := range oldRepositoriesList {
		if repository.Uuid == uuid {
			oldRepositoriesList[i].Status = status
			if status == StatusPulled || status == StatusPushed || status == StatusSynchronized {
				t := time.Now()
				oldRepositoriesList[i].UpdatedAt = t.Format(timeFormat)
			}
		}
	}

	err := configuration.SaveRepositoriesConfig(oldRepositoriesList)
	if err != nil {
		return
	}

	return
}

func UpdateRepositoryState(uuid string, state string) {

	oldRepositoriesList := configuration.GetRepositoriesConfig()
	for i, repository := range oldRepositoriesList {
		if repository.Uuid == uuid {
			oldRepositoriesList[i].State = state
			t := time.Now()
			oldRepositoriesList[i].UpdatedAt = t.Format(timeFormat)
		}
	}

	err := configuration.SaveRepositoriesConfig(oldRepositoriesList)
	if err != nil {
		return
	}

	return
}
