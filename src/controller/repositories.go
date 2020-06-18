package controller

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/component/gui"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/Shitovdm/GitRsync/src/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

const (
	// StatusInitiated describe status "initiated"
	StatusInitiated = "initiated"
	// StatusPendingPull describe status "pending_pull"
	StatusPendingPull = "pending_pull"
	// StatusPulled describe status "pulled"
	StatusPulled = "pulled"
	// StatusPullFailed describe status "pull_failed"
	StatusPullFailed = "pull_failed"
	// StatusPendingPush describe status "pending_push"
	StatusPendingPush = "pending_push"
	// StatusPushed describe status "pushed"
	StatusPushed = "pushed"
	// StatusPushFailed describe status "push_failed"
	StatusPushFailed = "push_failed"
	// StatusPendingClean describe status "pending_clear"
	StatusPendingClean = "pending_clear"
	// StatusCleaned describe status "cleared"
	StatusCleaned = "cleared"
	// StatusCleanFailed describe status "clear_failed"
	StatusCleanFailed = "clear_failed"
	// StatusSynchronized describe status "synced"
	StatusSynchronized = "synced"
	// StatusFailed describe status "failed"
	StatusFailed = "failed"
)

var (
	timeFormat = "02-01-2006 15:04"
)

// RepositoriesController struct describes repositories section controller.
type RepositoriesController struct{}

// Index describes repositories index page.
func (ctrl RepositoriesController) Index(c *gin.Context) {
	menu := gui.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Repositories"
	templateParams["repositories"], _ = conf.GetRepositoriesConfigData()
	templateParams["platforms"], _ = conf.GetPlatformsConfigData()
	c.HTML(http.StatusOK, "repositories/index", templateParams)
}

// Add describes add repository action.
func (ctrl RepositoriesController) Add(c *gin.Context) {

	var addRepositoryRequest model.AddRepositoryRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &addRepositoryRequest)
	if err != nil {
		logger.Error("RepositoriesController/Add", err.Error())
		return
	}

	newRepositoryUUID, _ := uuid.NewV4()
	repositories := conf.GetRepositoriesConfig()
	repositories = append(repositories, model.RepositoryConfig{
		UUID:                    newRepositoryUUID.String(),
		Name:                    addRepositoryRequest.Name,
		SourcePlatformUUID:      addRepositoryRequest.SourcePlatformUUID,
		SourcePlatformPath:      addRepositoryRequest.SourcePlatformPath,
		DestinationPlatformUUID: addRepositoryRequest.DestinationPlatformUUID,
		DestinationPlatformPath: addRepositoryRequest.DestinationPlatformPath,
		Status:                  StatusInitiated,
		State:                   "active",
		UpdatedAt:               "",
	})

	err = conf.SaveRepositoriesConfig(repositories)
	if err != nil {
		logger.Error("RepositoriesController/Add", err.Error())
	}

	logger.Info("RepositoriesController/Add", fmt.Sprintf("New repository with name %s added successfully!", addRepositoryRequest.Name))
}

// Edit describes edit repository action.
func (ctrl RepositoriesController) Edit(c *gin.Context) {

	var editRepositoryRequest model.EditRepositoryRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &editRepositoryRequest)
	if err != nil {
		logger.Error("RepositoriesController/Edit", err.Error())
		return
	}

	oldRepositoriesList := conf.GetRepositoriesConfig()
	newRepositoriesList := make([]model.RepositoryConfig, 0)
	for _, repository := range oldRepositoriesList {
		if repository.UUID == editRepositoryRequest.UUID {
			newRepositoriesList = append(newRepositoriesList, model.RepositoryConfig{
				UUID:                    repository.UUID,
				Name:                    editRepositoryRequest.Name,
				SourcePlatformUUID:      editRepositoryRequest.SourcePlatformUUID,
				SourcePlatformPath:      editRepositoryRequest.SourcePlatformPath,
				DestinationPlatformUUID: editRepositoryRequest.DestinationPlatformUUID,
				DestinationPlatformPath: editRepositoryRequest.DestinationPlatformPath,
				Status:                  repository.Status,
				State:                   repository.State,
				UpdatedAt:               repository.UpdatedAt,
			})
			continue
		}
		newRepositoriesList = append(newRepositoriesList, repository)
	}

	err = conf.SaveRepositoriesConfig(newRepositoriesList)
	if err != nil {
		logger.Error("RepositoriesController/Edit", err.Error())
	}

	logger.Info("RepositoriesController/Edit", fmt.Sprintf("Repository with name %s successfully edited!", editRepositoryRequest.Name))
}

// Remove describes remove repository action.
func (ctrl RepositoriesController) Remove(c *gin.Context) {

	var removeRepositoryRequest model.RemoveRepositoryRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &removeRepositoryRequest)
	if err != nil {
		logger.Error("RepositoriesController/Remove", err.Error())
		return
	}

	removedRepositoryName := ""
	oldRepositoriesList := conf.GetRepositoriesConfig()
	newRepositoriesList := make([]model.RepositoryConfig, 0)
	for _, repository := range oldRepositoriesList {
		if repository.UUID != removeRepositoryRequest.UUID {
			newRepositoriesList = append(newRepositoriesList, repository)
		} else {
			removedRepositoryName = repository.Name
		}
	}

	err = conf.SaveRepositoriesConfig(newRepositoriesList)
	if err != nil {
		logger.Error("RepositoriesController/Remove", err.Error())
	}

	logger.Info("RepositoriesController/Remove", fmt.Sprintf("Repository with name %s successfully removed!", removedRepositoryName))
}

// UpdateRepositoryStatus describes update repository status action.
func UpdateRepositoryStatus(uuid string, status string) {

	oldRepositoriesList := conf.GetRepositoriesConfig()
	for i, repository := range oldRepositoriesList {
		if repository.UUID == uuid {
			oldRepositoriesList[i].Status = status
			if status == StatusPulled || status == StatusPushed || status == StatusSynchronized {
				t := time.Now()
				oldRepositoriesList[i].UpdatedAt = t.Format(timeFormat)
			}
		}
	}

	err := conf.SaveRepositoriesConfig(oldRepositoriesList)
	if err != nil {
		return
	}
}

// UpdateRepositoryState describes update repository state action.
func UpdateRepositoryState(uuid string, state string) {

	oldRepositoriesList := conf.GetRepositoriesConfig()
	for i, repository := range oldRepositoriesList {
		if repository.UUID == uuid {
			oldRepositoriesList[i].State = state
			t := time.Now()
			oldRepositoriesList[i].UpdatedAt = t.Format(timeFormat)
		}
	}

	err := conf.SaveRepositoriesConfig(oldRepositoriesList)
	if err != nil {
		return
	}
}
