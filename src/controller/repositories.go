package controller

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/gui"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/Shitovdm/GitRsync/src/model/platform"
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
	templateParams["repositories"], _ = repository.GetAllInInterface()
	templateParams["platforms"], _ = platform.GetAllInInterface()
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

	repo := repository.Get(editRepositoryRequest.UUID)
	repo.SetName(editRepositoryRequest.Name)
	repo.SetSourcePlatformPath(editRepositoryRequest.SourcePlatformPath)
	repo.SetSourcePlatformUUID(editRepositoryRequest.SourcePlatformUUID)
	repo.SetDestinationPlatformPath(editRepositoryRequest.DestinationPlatformPath)
	repo.SetDestinationPlatformUUID(editRepositoryRequest.DestinationPlatformUUID)
	err = repo.Update()
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

	repo := repository.Get(removeRepositoryRequest.UUID)
	err = repo.Delete()
	if err != nil {
		logger.Error("RepositoriesController/Remove", err.Error())
	}

	logger.Info("RepositoriesController/Remove", fmt.Sprintf("Repository with name %s successfully removed!", repo.GetName()))
}
