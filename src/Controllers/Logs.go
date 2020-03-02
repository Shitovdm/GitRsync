package Controllers

import (
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LogsController struct{}

func (ctrl LogsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	c.HTML(http.StatusOK, "logs/index", templateParams)
}