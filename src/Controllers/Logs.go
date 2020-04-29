package Controllers

import (
	"fmt"
	"github.com/Shitovdm/git-rsync/src/Components/Helpers"
	"github.com/Shitovdm/git-rsync/src/Components/Interface"
	"github.com/Shitovdm/git-rsync/src/Components/Logger"
	"github.com/Shitovdm/git-rsync/src/Models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type LogsController struct{}

func (ctrl LogsController) Index(c *gin.Context) {
	menu := Interface.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Logs"
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

func (ctrl LogsController) RemoveRuntime(c *gin.Context) {

	var removeRuntimeLogsRequest Models.RuntimeLogsRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &removeRuntimeLogsRequest)
	if err != nil {
		Msg := fmt.Sprintf("Error while removing runtime logs! %s", err.Error())
		Logger.Error("LogsController/RemoveRuntime", Msg)
		return
	}

	err = Logger.ClearRuntimeLogs()
	if err != nil {
		Msg := fmt.Sprintf("%s", err.Error())
		Logger.Error("LogsController/RemoveRuntime", Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(Msg)))
		return
	}

	Msg := "Runtime logs successfully removed!"
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg)))
	Logger.Info("LogsController/RemoveRuntime", Msg)
	return
}

func (ctrl LogsController) RemoveAll(c *gin.Context) {

	var removeAllLogsRequest Models.RuntimeLogsRequest
	conn, err := Helpers.WsHandler(c.Writer, c.Request, &removeAllLogsRequest)
	if err != nil {
		Msg := fmt.Sprintf("Error while removing all logs! %s", err.Error())
		Logger.Error("LogsController/RemoveAll", Msg)
		return
	}

	err = Logger.ClearAllLogs()
	if err != nil {
		Msg := fmt.Sprintf("%s", err.Error())
		Logger.Error("LogsController/RemoveAll", Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(Msg)))
		return
	}

	Msg := "All logs successfully removed!"
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg)))
	Logger.Info("LogsController/RemoveAll", Msg)
	return
}
