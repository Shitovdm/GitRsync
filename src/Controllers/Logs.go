package Controllers

import (
	"fmt"
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
		_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\x1b[92m%s\x1b[0m", " GitRsync (C) Shitov Dmitry")))
		for _, logNote := range Logger.GetRuntimeLogs() {
			//	Only fot current session.
			if logNote.SessionID == Logger.GetSessionId() {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(Logger.BuildRuntimeLogNote(logNote)))
			}
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
