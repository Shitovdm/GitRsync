package Controllers

import (
	"fmt"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Logger"
	"github.com/Shitovdm/git-repo-exporter/src/Models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type LogsController struct{}

func (ctrl LogsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}

	Logger.GetRuntimeLogFile()
	//fmt.Println(Logger.GetRuntimeLogs())

	Logger.Info("LogsController", "Log file successfully loaded!")

	c.HTML(http.StatusOK, "logs/index", templateParams)
}

func (ctrl LogsController) Process(c *gin.Context) {




	var processLogsRequest Models.ProcessLogsRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &processLogsRequest)
	if err != nil {
		log.Println(err.Error())
		return
	}

	switch processLogsRequest.Action {
	case "init":
		for _, logNote := range Logger.GetRuntimeLogs() {
			runtimeLog := "[" + logNote.Time + "]" + "\t"
			runtimeLog += logNote.SessionID + "\t"
			runtimeLog += logNote.Level + "\t"
			runtimeLog += logNote.Category + "\t"
			runtimeLog += logNote.Message
			_ = conn.WriteMessage(websocket.TextMessage, []byte(Logger.SetLogLevel(logNote.Level, runtimeLog)))
		}


		//c.JSON(http.StatusOK, `{"data": "ret"}`)
	}

	fmt.Println("connected to append method!")
}