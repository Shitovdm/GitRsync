package controller

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/component/gui"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/Shitovdm/GitRsync/src/model"
	"github.com/Shitovdm/GitRsync/src/model/repository"
	"github.com/gin-gonic/gin"
	"net/http"
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

	var addRepositoryRequest repository.AddRepositoryRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &addRepositoryRequest)
	if err != nil {
		logger.Error("RepositoriesController/Add", err.Error())
		return
	}

	err = repository.Create(&repository.Repository{
		Name:                    addRepositoryRequest.Name,
		SourcePlatformUUID:      addRepositoryRequest.SourcePlatformUUID,
		SourcePlatformPath:      addRepositoryRequest.SourcePlatformPath,
		DestinationPlatformUUID: addRepositoryRequest.DestinationPlatformUUID,
		DestinationPlatformPath: addRepositoryRequest.DestinationPlatformPath,
	})
	if err != nil {
		logger.Error("RepositoriesController/Add", err.Error())
		return
	}

	logger.Info("RepositoriesController/Add", fmt.Sprintf("New repository with name %s added successfully!", addRepositoryRequest.Name))
}

// Edit describes edit repository action.
func (ctrl RepositoriesController) Edit(c *gin.Context) {

	var editRepositoryRequest repository.EditRepositoryRequest
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

	var removeRepositoryRequest repository.RemoveRepositoryRequest
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
