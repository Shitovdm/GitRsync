package controller

import (
	"encoding/json"
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/component/gui"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/Shitovdm/GitRsync/src/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// SettingsController struct describes settings section controller.
type SettingsController struct{}

// Index describes settings index page.
func (ctrl SettingsController) Index(c *gin.Context) {
	menu := gui.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Settings"
	templateParams["appconfig"], _ = conf.GetAppConfigData()
	c.HTML(http.StatusOK, "settings/index", templateParams)
}

// Save describes save settings action.
func (ctrl SettingsController) Save(c *gin.Context) {
	var saveSettingsRequest model.SaveSettingsRequest
	_, err := helper.WsHandler(c.Writer, c.Request, &saveSettingsRequest)
	if err != nil {
		logger.Error("SettingsController/Save", err.Error())
		return
	}

	appConfig := conf.GetAppConfig()
	section := saveSettingsRequest.Section
	field := saveSettingsRequest.Field

	reflectValueTypeNeeded := reflect.Indirect(reflect.ValueOf(appConfig)).FieldByName(section).FieldByName(field).Type()
	reflectValue := reflect.ValueOf(saveSettingsRequest.Value)
	switch reflectValueTypeNeeded.String() {
	case "string":

	case "int":
		nonFractionalPart := strings.Split(saveSettingsRequest.Value.(string), ".")
		val, _ := strconv.Atoi(nonFractionalPart[0])
		reflectValue = reflect.ValueOf(val)
	case "bool":
		if saveSettingsRequest.Value == "true" {
			reflectValue = reflect.ValueOf(true)
		}
		if saveSettingsRequest.Value == "false" {
			reflectValue = reflect.ValueOf(false)
		}
	default:
		//	structs
		byteData, _ := json.Marshal(saveSettingsRequest.Value)
		switch reflectValueTypeNeeded.String() {
		case "[]model.CommittersRule":
			val := make([]model.CommittersRule, 0)
			_ = json.Unmarshal(byteData, &val)
			reflectValue = reflect.ValueOf(val)
		case "model.GitUser":
			val := model.GitUser{}
			_ = json.Unmarshal(byteData, &val)
			reflectValue = reflect.ValueOf(val)
		}
	}

	reflect.Indirect(reflect.ValueOf(appConfig)).FieldByName(section).FieldByName(field).Set(reflectValue)
	err = conf.SaveAppConfig(appConfig)
	if err != nil {
		logger.Error("SettingsController/Save", err.Error())
		return
	}
}

// OpenRawConfig opens raw config file.
func (ctrl SettingsController) OpenRawConfig(_ *gin.Context) {
	helper.ExploreDir(conf.BuildPlatformPath(``))
}

// ExploreConfigDir opens config dir in explorer.
func (ctrl SettingsController) ExploreConfigDir(_ *gin.Context) {
	helper.ExploreDir(conf.BuildPlatformPath(`projects`))
}