package Controllers

import (
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Logger"
	"github.com/Shitovdm/git-repo-exporter/src/Models"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

type PlatformsController struct{}

func (ctrl PlatformsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Platforms"
	templateParams["platforms"], _ = Configuration.GetPlatformsConfigData()
	c.HTML(http.StatusOK, "platforms/index", templateParams)
}

func (ctrl PlatformsController) Add(c *gin.Context) {

	var addPlatformRequest Models.AddPlatformRequest
	_, err := Helpers.WsHandler(c.Writer, c.Request, &addPlatformRequest)
	if err != nil {
		Logger.Error("PlatformsController/Add", err.Error())
		return
	}

	newPlatformUuid, _ := uuid.NewV4()
	platforms := Configuration.GetPlatformsConfig()
	platforms = append(platforms, Models.PlatformConfig{
		Uuid:     newPlatformUuid.String(),
		Name:     addPlatformRequest.Name,
		Address:  addPlatformRequest.Address,
		Username: addPlatformRequest.Username,
		Password: addPlatformRequest.Password,
	})

	err = Configuration.SavePlatformsConfig(platforms)
	if err != nil {
		Logger.Error("PlatformsController/Add", err.Error())
	}

	return
}

func (ctrl PlatformsController) Edit(c *gin.Context) {

	var editPlatformRequest Models.EditPlatformRequest
	_, err := Helpers.WsHandler(c.Writer, c.Request, &editPlatformRequest)
	if err != nil {
		Logger.Error("PlatformsController/Edit", err.Error())
		return
	}

	oldPlatformsList := Configuration.GetPlatformsConfig()
	newPlatformsList := make([]Models.PlatformConfig, 0)
	for _, platform := range oldPlatformsList {
		if platform.Uuid == editPlatformRequest.Uuid {
			newPlatformsList = append(newPlatformsList, Models.PlatformConfig{
				Uuid:     platform.Uuid,
				Name:     editPlatformRequest.Name,
				Address:  editPlatformRequest.Address,
				Username: editPlatformRequest.Username,
				Password: editPlatformRequest.Password,
			})
			continue
		}
		newPlatformsList = append(newPlatformsList, platform)
	}

	err = Configuration.SavePlatformsConfig(newPlatformsList)
	if err != nil {
		Logger.Error("PlatformsController/Edit", err.Error())
	}

	return
}

func (ctrl PlatformsController) Remove(c *gin.Context) {

	var removePlatformRequest Models.RemovePlatformRequest
	_, err := Helpers.WsHandler(c.Writer, c.Request, &removePlatformRequest)
	if err != nil {
		Logger.Error("PlatformsController/Remove", err.Error())
		return
	}

	oldPlatformsList := Configuration.GetPlatformsConfig()
	newPlatformsList := make([]Models.PlatformConfig, 0)
	for _, platform := range oldPlatformsList {
		if platform.Uuid != removePlatformRequest.Uuid {
			newPlatformsList = append(newPlatformsList, platform)
		}
	}

	err = Configuration.SavePlatformsConfig(newPlatformsList)
	if err != nil {
		Logger.Error("PlatformsController/Remove", err.Error())
	}

	return
}
