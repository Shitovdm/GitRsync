package Controllers

import (
	"github.com/Shitovdm/git-repo-exporter/src/Components/Configuration"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexController struct{}

func (ctrl IndexController) Index(c *gin.Context) {

	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Dashboard"
	templateParams["config"] = Configuration.GetAppConfig()
	templateParams["platforms"] = Configuration.GetPlatformsConfig()
	templateParams["repositories"] = Configuration.GetRepositoriesConfig()

	c.HTML(http.StatusOK, "index/index", templateParams)
}
