package controllers

import (
	"github.com/Shitovdm/GitRsync/src/components/configuration"
	"github.com/Shitovdm/GitRsync/src/components/interface"
	"github.com/Shitovdm/GitRsync/src/components/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexController struct{}

func (ctrl IndexController) Index(c *gin.Context) {

	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Dashboard"
	templateParams["config"], _ = configuration.GetAppConfigData()
	templateParams["platforms"], _ = configuration.GetPlatformsConfigData()
	templateParams["active_repositories"], _ = configuration.GetActiveRepositoriesConfigData()
	templateParams["blocked_repositories"], _ = configuration.GetBlockedRepositoriesConfigData()
	templateParams["log_error_count"] = logger.CountErrorsInRuntimeLog()

	c.HTML(http.StatusOK, "index/index", templateParams)
}
