package Controllers

import (
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PlatformsController struct{}

func (ctrl PlatformsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Platforms"

	_, conferr := Helpers.GetAppConfig()
	if conferr != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	c.HTML(http.StatusOK, "platforms/index", templateParams)
}
