package Controllers

import (
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexController struct{}

func (ctrl IndexController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Dashboard"
	templateParams["config"] = Helpers.GetAppConfig()
	templateParams["platforms"] = Helpers.GetPlatformsConfig()
	templateParams["repositories"] = Helpers.GetRepositoriesConfig()

	c.HTML(http.StatusOK, "index/index", templateParams)
}
