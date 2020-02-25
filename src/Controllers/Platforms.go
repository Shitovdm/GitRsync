package Controllers

import (
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/Shitovdm/git-repo-exporter/src/Models"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
)

type PlatformsController struct{}

func (ctrl PlatformsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Platforms"
	templateParams["platforms"], _ = Helpers.GetPlatformsConfigData()

	c.HTML(http.StatusOK, "platforms/index", templateParams)
}

func (ctrl PlatformsController) Add(c *gin.Context) {

	var addPlatformRequest Models.AddPlatformRequest
	err := Helpers.WsHandler(c.Writer, c.Request, &addPlatformRequest)
	if err != nil {
		log.Println(err.Error())
		return
	}

	newPlatformUuid, _ := uuid.NewV4()
	platforms := Helpers.GetPlatformsConfig()
	platforms = append(platforms, Models.PlatformConfig{
		Guid:     newPlatformUuid.String(),
		Name:     addPlatformRequest.Name,
		Address:  addPlatformRequest.Address,
		Username: addPlatformRequest.Username,
		Password: addPlatformRequest.Password,
	})

	err = Helpers.SavePlatformsConfig(platforms)
	if err != nil {
		log.Println(err.Error())
	}

	return
}

func (ctrl PlatformsController) Edit(c *gin.Context) {

	var editPlatformRequest Models.EditPlatformRequest
	err := Helpers.WsHandler(c.Writer, c.Request, &editPlatformRequest)
	if err != nil {
		log.Println(err.Error())
		return
	}

	oldPlatformsList := Helpers.GetPlatformsConfig()
	newPlatformsList := make([]Models.PlatformConfig, 0)
	for _, platform := range oldPlatformsList {
		if platform.Guid == editPlatformRequest.Guid {
			newPlatformsList = append(newPlatformsList, Models.PlatformConfig{
				Guid:     platform.Guid,
				Name:     editPlatformRequest.Name,
				Address:  editPlatformRequest.Address,
				Username: editPlatformRequest.Username,
				Password: editPlatformRequest.Password,
			})
			continue
		}
		newPlatformsList = append(newPlatformsList, platform)
	}

	err = Helpers.SavePlatformsConfig(newPlatformsList)
	if err != nil {
		log.Println(err.Error())
	}

	return
}

func (ctrl PlatformsController) Remove(c *gin.Context) {

	var removePlatformRequest Models.RemovePlatformRequest
	err := Helpers.WsHandler(c.Writer, c.Request, &removePlatformRequest)
	if err != nil {
		log.Println(err.Error())
		return
	}

	oldPlatformsList := Helpers.GetPlatformsConfig()
	newPlatformsList := make([]Models.PlatformConfig, 0)
	for _, platform := range oldPlatformsList {
		if platform.Guid != removePlatformRequest.Guid {
			newPlatformsList = append(newPlatformsList, platform)
		}
	}

	err = Helpers.SavePlatformsConfig(newPlatformsList)
	if err != nil {
		log.Println(err.Error())
	}

	return
}
