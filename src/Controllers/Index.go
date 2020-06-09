package Controllers

import (
	"github.com/Shitovdm/git-rsync/src/Components/Configuration"
	"github.com/Shitovdm/git-rsync/src/Components/Interface"
	"github.com/Shitovdm/git-rsync/src/Components/Logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexController struct{}

func (ctrl IndexController) Index(c *gin.Context) {

	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Dashboard"
	templateParams["config"], _ = Configuration.GetAppConfigData()
	templateParams["platforms"], _ = Configuration.GetPlatformsConfigData()
	templateParams["active_repositories"], _ = Configuration.GetActiveRepositoriesConfigData()
	templateParams["blocked_repositories"], _ = Configuration.GetBlockedRepositoriesConfigData()
	templateParams["log_error_count"] = Logger.CountErrorsInRuntimeLog()

	c.HTML(http.StatusOK, "index/index", templateParams)
}
