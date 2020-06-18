package controllers

import (
	"encoding/json"
	"github.com/Shitovdm/GitRsync/src/components/configuration"
	"github.com/Shitovdm/GitRsync/src/components/helpers"
	"github.com/Shitovdm/GitRsync/src/components/interface"
	"github.com/Shitovdm/GitRsync/src/components/logger"
	"github.com/Shitovdm/GitRsync/src/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type SettingsController struct{}

func (ctrl SettingsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Settings"
	templateParams["appconfig"], _ = configuration.GetAppConfigData()
	c.HTML(http.StatusOK, "settings/index", templateParams)
}

func (ctrl SettingsController) Save(c *gin.Context) {
	var saveSettingsRequest models.SaveSettingsRequest
	_, err := helpers.WsHandler(c.Writer, c.Request, &saveSettingsRequest)
	if err != nil {
		logger.Error("SettingsController/Save", err.Error())
		return
	}

	appConfig := configuration.GetAppConfig()
	section := saveSettingsRequest.Section
	field := saveSettingsRequest.Field

	reflectValueTypeNeeded := reflect.Indirect(reflect.ValueOf(appConfig)).FieldByName(section).FieldByName(field).Type()
	reflectValue := reflect.ValueOf(saveSettingsRequest.Value)
	switch reflectValueTypeNeeded.String() {
	case "string":

		break
	case "int":
		nonFractionalPart := strings.Split(saveSettingsRequest.Value.(string), ".")
		val, _ := strconv.Atoi(nonFractionalPart[0])
		reflectValue = reflect.ValueOf(val)
		break
	case "bool":
		if saveSettingsRequest.Value == "true" {
			reflectValue = reflect.ValueOf(true)
		}
		if saveSettingsRequest.Value == "false" {
			reflectValue = reflect.ValueOf(false)
		}
		break
	default:
		//	structs
		byteData, _ := json.Marshal(saveSettingsRequest.Value)
		switch reflectValueTypeNeeded.String() {
		case "[]Models.CommittersRule":
			val := make([]models.CommittersRule, 0)
			_ = json.Unmarshal(byteData, &val)
			reflectValue = reflect.ValueOf(val)
			break
		case "Models.GitUser":
			val := models.GitUser{}
			_ = json.Unmarshal(byteData, &val)
			reflectValue = reflect.ValueOf(val)
			break
		}
		break
	}

	reflect.Indirect(reflect.ValueOf(appConfig)).FieldByName(section).FieldByName(field).Set(reflectValue)
	err = configuration.SaveAppConfig(appConfig)
	if err != nil {
		logger.Error("SettingsController/Save", err.Error())
		return
	}

	return
}
