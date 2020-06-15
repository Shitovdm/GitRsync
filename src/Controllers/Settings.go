package Controllers

import (
	"github.com/Shitovdm/git-rsync/src/Components/Configuration"
	"github.com/Shitovdm/git-rsync/src/Components/Helpers"
	"github.com/Shitovdm/git-rsync/src/Components/Interface"
	"github.com/Shitovdm/git-rsync/src/Components/Logger"
	"github.com/Shitovdm/git-rsync/src/Models"
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
	templateParams["appconfig"], _ = Configuration.GetAppConfigData()
	c.HTML(http.StatusOK, "settings/index", templateParams)
}

func (ctrl SettingsController) Save(c *gin.Context) {

	var saveSettingsRequest Models.SaveSettingsRequest
	_, err := Helpers.WsHandler(c.Writer, c.Request, &saveSettingsRequest)
	if err != nil {
		Logger.Error("SettingsController/Save", err.Error())
		return
	}

	appConfig := Configuration.GetAppConfig()
	section := saveSettingsRequest.Section
	field := saveSettingsRequest.Field

	reflectValueTypeNeeded := reflect.Indirect(reflect.ValueOf(appConfig)).FieldByName(section).FieldByName(field).Type().String()
	reflectValue := reflect.ValueOf(saveSettingsRequest.Value)

	if reflectValueTypeNeeded == "int" {
		nonFractionalPart := strings.Split(saveSettingsRequest.Value.(string), ".")
		val, _ := strconv.Atoi(nonFractionalPart[0])
		reflectValue = reflect.ValueOf(val)
	}

	if reflectValueTypeNeeded == "bool" {
		if saveSettingsRequest.Value == "true" {
			reflectValue = reflect.ValueOf(true)
		}
		if saveSettingsRequest.Value == "false" {
			reflectValue = reflect.ValueOf(false)
		}
	}

	reflect.Indirect(reflect.ValueOf(appConfig)).FieldByName(section).FieldByName(field).Set(reflectValue)
	err = Configuration.SaveAppConfig(appConfig)
	if err != nil {
		Logger.Error("SettingsController/Save", err.Error())
		return
	}

	return
}
