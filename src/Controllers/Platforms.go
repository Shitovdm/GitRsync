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
)

type PlatformsController struct{}

func (ctrl PlatformsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Platforms"
	templateParams["platforms"], _ = configuration.GetPlatformsConfigData()
	c.HTML(http.StatusOK, "platforms/index", templateParams)
}

func (ctrl PlatformsController) Add(c *gin.Context) {

	var addPlatformRequest models.AddPlatformRequest
	_, err := helpers.WsHandler(c.Writer, c.Request, &addPlatformRequest)
	if err != nil {
		logger.Error("PlatformsController/Add", err.Error())
		return
	}

	newPlatformUuid, _ := uuid.NewV4()
	platforms := configuration.GetPlatformsConfig()
	platforms = append(platforms, models.PlatformConfig{
		Uuid:     newPlatformUuid.String(),
		Name:     addPlatformRequest.Name,
		Address:  addPlatformRequest.Address,
		Username: addPlatformRequest.Username,
		Password: addPlatformRequest.Password,
	})

	err = configuration.SavePlatformsConfig(platforms)
	if err != nil {
		logger.Error("PlatformsController/Add", err.Error())
	}

	logger.Info("PlatformsController/Add", fmt.Sprintf("New platform with name %s added successfully!", addPlatformRequest.Name))
	return
}

func (ctrl PlatformsController) Edit(c *gin.Context) {

	var editPlatformRequest models.EditPlatformRequest
	_, err := helpers.WsHandler(c.Writer, c.Request, &editPlatformRequest)
	if err != nil {
		logger.Error("PlatformsController/Edit", err.Error())
		return
	}

	oldPlatformsList := configuration.GetPlatformsConfig()
	newPlatformsList := make([]models.PlatformConfig, 0)
	for _, platform := range oldPlatformsList {
		if platform.Uuid == editPlatformRequest.Uuid {
			newPlatformsList = append(newPlatformsList, models.PlatformConfig{
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

	err = configuration.SavePlatformsConfig(newPlatformsList)
	if err != nil {
		logger.Error("PlatformsController/Edit", err.Error())
	}

	logger.Info("PlatformsController/Edit", fmt.Sprintf("Platform with name %s successfully edited!", editPlatformRequest.Name))
	return
}

func (ctrl PlatformsController) Remove(c *gin.Context) {

	var removePlatformRequest models.RemovePlatformRequest
	_, err := helpers.WsHandler(c.Writer, c.Request, &removePlatformRequest)
	if err != nil {
		logger.Error("PlatformsController/Remove", err.Error())
		return
	}

	removedPlatformName := ""
	oldPlatformsList := configuration.GetPlatformsConfig()
	newPlatformsList := make([]models.PlatformConfig, 0)
	for _, platform := range oldPlatformsList {
		if platform.Uuid != removePlatformRequest.Uuid {
			newPlatformsList = append(newPlatformsList, platform)
		} else {
			removedPlatformName = platform.Name
		}
	}

	err = configuration.SavePlatformsConfig(newPlatformsList)
	if err != nil {
		logger.Error("PlatformsController/Remove", err.Error())
	}

	logger.Info("PlatformsController/Remove", fmt.Sprintf("Platform with name %s successfully removed!", removedPlatformName))
	return
}
