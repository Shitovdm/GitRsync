package Controllers

import (
	"github.com/Shitovdm/git-repo-exporter/src/Components/Helpers"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Interface"
	"github.com/Shitovdm/git-repo-exporter/src/Components/Logger"
	"github.com/Shitovdm/git-repo-exporter/src/Models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type LogsController struct{}

func (ctrl LogsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}

	//Logger.GetRuntimeLogFile()
	//fmt.Println(Logger.GetRuntimeLogs())
	//Logger.Info("LogsController", "Log file successfully loaded!")

	c.HTML(http.StatusOK, "logs/index", templateParams)
}

func (ctrl LogsController) Subscribe(c *gin.Context) {
	var subsctibeToLogRequest Models.RuntimeLogsRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &subsctibeToLogRequest)
	if err != nil {
		Logger.Error("LogsController/Subscribe", err.Error())
		return
	}

	switch subsctibeToLogRequest.Action {
	case "init":
		_ = conn.WriteMessage(websocket.TextMessage, []byte(" GitRsync (C) Shitov Dmitry"))
		for _, logNote := range Logger.GetRuntimeLogs() {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(Logger.BuildRuntimeLogNote(logNote)))
		}
		go func() {
			for {
				if Logger.GetRuntimeLogNote() != "" {
					_ = conn.WriteMessage(websocket.TextMessage, []byte(Logger.GetRuntimeLogNote()))
					Logger.ResetRuntimeLogNote()
				}
				time.Sleep(500 * time.Microsecond)
			}
		}()
	}
}