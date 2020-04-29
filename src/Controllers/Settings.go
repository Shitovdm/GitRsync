package Controllers

import (
	"github.com/Shitovdm/git-rsync/src/Components/Interface"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SettingsController struct{}

func (ctrl SettingsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Settings"
	c.HTML(http.StatusOK, "settings/index", templateParams)
}
