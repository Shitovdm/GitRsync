package controller

import (
	"github.com/Shitovdm/GitRsync/src/component/conf"
	"github.com/Shitovdm/GitRsync/src/component/gui"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// IndexController struct describes index section controller.
type IndexController struct{}

// Index describes index index page.
func (ctrl IndexController) Index(c *gin.Context) {

	menu := gui.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Dashboard"
	templateParams["config"], _ = conf.GetAppConfigData()
	templateParams["platforms"], _ = conf.GetPlatformsConfigData()
	templateParams["active_repositories"], _ = conf.GetActiveRepositoriesConfigData()
	templateParams["blocked_repositories"], _ = conf.GetBlockedRepositoriesConfigData()
	templateParams["log_error_count"] = logger.CountErrorsInRuntimeLog()

	c.HTML(http.StatusOK, "index/index", templateParams)
}
