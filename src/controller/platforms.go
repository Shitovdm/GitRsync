package controller

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/gui"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/Shitovdm/GitRsync/src/model/platform"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PlatformsController struct describes platforms section controller.
type PlatformsController struct{}

// Index describes platforms index page.
func (ctrl PlatformsController) Index(c *gin.Context) {

	menu := gui.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Platforms"
	templateParams["platforms"], _ = platform.GetAllInInterface()

	c.HTML(http.StatusOK, "platforms/index", templateParams)
}

// Add describes add platform action.
func (ctrl PlatformsController) Add(c *gin.Context) {

	var addPlatformRequest platform.AddPlatformRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &addPlatformRequest)
	if err != nil {
		logger.Error("PlatformsController/Add", err.Error())
		return
	}

	err = platform.Create(&platform.Platform{
		Name:     addPlatformRequest.Name,
		Address:  addPlatformRequest.Address,
		Username: addPlatformRequest.Username,
		Password: addPlatformRequest.Password,
	})
	if err != nil {
		logger.Error("PlatformsController/Add", err.Error())
	}

	logger.Info("PlatformsController/Add", fmt.Sprintf("New platform with name %s added successfully!", addPlatformRequest.Name))
}

// Edit describes edit platform action.
func (ctrl PlatformsController) Edit(c *gin.Context) {

	var editPlatformRequest platform.EditPlatformRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &editPlatformRequest)
	if err != nil {
		logger.Error("PlatformsController/Edit", err.Error())
		return
	}

	pl := platform.Get(editPlatformRequest.UUID)
	pl.SetName(editPlatformRequest.Name)
	pl.SetAddress(editPlatformRequest.Address)
	pl.SetUsername(editPlatformRequest.Username)
	pl.SetPassword(editPlatformRequest.Password)
	err = pl.Update()
	if err != nil {
		logger.Error("PlatformsController/Edit", err.Error())
	}

	logger.Info("PlatformsController/Edit", fmt.Sprintf("Platform with name %s successfully edited!", editPlatformRequest.Name))
}

// Remove describes remove platform action.
func (ctrl PlatformsController) Remove(c *gin.Context) {

	var removePlatformRequest platform.RemovePlatformRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &removePlatformRequest)
	if err != nil {
		logger.Error("PlatformsController/Remove", err.Error())
		return
	}

	pl := platform.Get(removePlatformRequest.UUID)
	err = pl.Delete()
	if err != nil {
		logger.Error("PlatformsController/Remove", err.Error())
	}

	logger.Info("PlatformsController/Remove", fmt.Sprintf("Platform with name %s successfully removed!", pl.GetName()))
}
