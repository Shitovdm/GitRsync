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
)

type PlatformsController struct{}

func (ctrl PlatformsController) Index(c *gin.Context) {
	menu := gui.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Platforms"
	templateParams["platforms"], _ = conf.GetPlatformsConfigData()
	c.HTML(http.StatusOK, "platforms/index", templateParams)
}

func (ctrl PlatformsController) Add(c *gin.Context) {

	var addPlatformRequest model.AddPlatformRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &addPlatformRequest)
	if err != nil {
		logger.Error("PlatformsController/Add", err.Error())
		return
	}

	newPlatformUuid, _ := uuid.NewV4()
	platforms := conf.GetPlatformsConfig()
	platforms = append(platforms, model.PlatformConfig{
		UUID:     newPlatformUuid.String(),
		Name:     addPlatformRequest.Name,
		Address:  addPlatformRequest.Address,
		Username: addPlatformRequest.Username,
		Password: addPlatformRequest.Password,
	})

	err = conf.SavePlatformsConfig(platforms)
	if err != nil {
		logger.Error("PlatformsController/Add", err.Error())
	}

	logger.Info("PlatformsController/Add", fmt.Sprintf("New platform with name %s added successfully!", addPlatformRequest.Name))
	return
}

func (ctrl PlatformsController) Edit(c *gin.Context) {

	var editPlatformRequest model.EditPlatformRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &editPlatformRequest)
	if err != nil {
		logger.Error("PlatformsController/Edit", err.Error())
		return
	}

	oldPlatformsList := conf.GetPlatformsConfig()
	newPlatformsList := make([]model.PlatformConfig, 0)
	for _, platform := range oldPlatformsList {
		if platform.UUID == editPlatformRequest.UUID {
			newPlatformsList = append(newPlatformsList, model.PlatformConfig{
				UUID:     platform.UUID,
				Name:     editPlatformRequest.Name,
				Address:  editPlatformRequest.Address,
				Username: editPlatformRequest.Username,
				Password: editPlatformRequest.Password,
			})
			continue
		}
		newPlatformsList = append(newPlatformsList, platform)
	}

	err = conf.SavePlatformsConfig(newPlatformsList)
	if err != nil {
		logger.Error("PlatformsController/Edit", err.Error())
	}

	logger.Info("PlatformsController/Edit", fmt.Sprintf("Platform with name %s successfully edited!", editPlatformRequest.Name))
	return
}

func (ctrl PlatformsController) Remove(c *gin.Context) {

	var removePlatformRequest model.RemovePlatformRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &removePlatformRequest)
	if err != nil {
		logger.Error("PlatformsController/Remove", err.Error())
		return
	}

	removedPlatformName := ""
	oldPlatformsList := conf.GetPlatformsConfig()
	newPlatformsList := make([]model.PlatformConfig, 0)
	for _, platform := range oldPlatformsList {
		if platform.UUID != removePlatformRequest.UUID {
			newPlatformsList = append(newPlatformsList, platform)
		} else {
			removedPlatformName = platform.Name
		}
	}

	err = conf.SavePlatformsConfig(newPlatformsList)
	if err != nil {
		logger.Error("PlatformsController/Remove", err.Error())
	}

	logger.Info("PlatformsController/Remove", fmt.Sprintf("Platform with name %s successfully removed!", removedPlatformName))
	return
}
