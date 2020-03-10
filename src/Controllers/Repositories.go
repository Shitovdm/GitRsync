package Controllers

import (
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Logger"
	"github.com/Shitovdm/git-repo-exporter/src/Models"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

type RepositoriesController struct{}

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
