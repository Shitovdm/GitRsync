package Controllers

import (
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LogsController struct{}

func (ctrl LogsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}

	Logger.GetRuntimeLogFile()
	fmt.Println(Logger.GetRuntimeLogs())

	Logger.Info("LogsController", "Log file successfully loaded!")

	c.HTML(http.StatusOK, "logs/index", templateParams)
}

