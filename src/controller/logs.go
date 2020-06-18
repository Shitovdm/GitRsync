package controller

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/component/gui"
	"github.com/Shitovdm/GitRsync/src/component/helper"
	"github.com/Shitovdm/GitRsync/src/component/logger"
	"github.com/Shitovdm/GitRsync/src/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// LogsController struct describes logs section controller.
type LogsController struct{}

// Index describes logs index page.
func (ctrl LogsController) Index(c *gin.Context) {

	menu := gui.GetMenu(c)
	templateParams := gin.H{"menu": menu}
	templateParams["title"] = "Logs"

	c.HTML(http.StatusOK, "logs/index", templateParams)
}

// Subscribe describes subscribe to log action.
func (ctrl LogsController) Subscribe(c *gin.Context) {

	var subsctibeToLogRequest model.RuntimeLogsRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &subsctibeToLogRequest)
	if err != nil {
		logger.Error("LogsController/Subscribe", err.Error())
		return
	}

	switch subsctibeToLogRequest.Action {
	case "init":
		_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\x1b[92m%s\x1b[0m", " GitRsync (C) Shitov Dmitry")))
		for _, logNote := range logger.GetRuntimeLogs() {
			//	Only fot current session.
			if logNote.SessionID == logger.GetSessionID() {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(logger.BuildRuntimeLogNote(logNote)))
			}
		}
		go func() {
			for {
				if logger.GetRuntimeLogNote() != "" {
					_ = conn.WriteMessage(websocket.TextMessage, []byte(logger.GetRuntimeLogNote()))
					logger.ResetRuntimeLogNote()
				}
				time.Sleep(500 * time.Microsecond)
			}
		}()
	}
}

// RemoveRuntime describes remove runtime logs action.
func (ctrl LogsController) RemoveRuntime(c *gin.Context) {

	var removeRuntimeLogsRequest model.RuntimeLogsRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &removeRuntimeLogsRequest)
	if err != nil {
		Msg := fmt.Sprintf("Error while removing runtime logs! %s", err.Error())
		logger.Error("LogsController/RemoveRuntime", Msg)
		return
	}

	err = logger.ClearRuntimeLogs()
	if err != nil {
		Msg := fmt.Sprintf("Error while cleaning runtime logs! %s", err.Error())
		logger.Error("LogsController/RemoveRuntime", Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(Msg)))
		return
	}

	Msg := "Runtime logs successfully removed!"
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONSuccess(Msg, nil)))
	logger.Info("LogsController/RemoveRuntime", Msg)
}

// RemoveAll describes remove all logs action.
func (ctrl LogsController) RemoveAll(c *gin.Context) {

	var removeAllLogsRequest model.RuntimeLogsRequest
	conn, err := helper.WsHandler(c.Writer, c.Request, &removeAllLogsRequest)
	if err != nil {
		Msg := fmt.Sprintf("Error while removing all logs! %s", err.Error())
		logger.Error("LogsController/RemoveAll", Msg)
		return
	}

	err = logger.ClearAllLogs()
	if err != nil {
		Msg := fmt.Sprintf("Error while cleaning all logs! %s", err.Error())
		logger.Error("LogsController/RemoveAll", Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONError(Msg)))
		return
	}

	Msg := "All logs successfully removed!"
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJSONSuccess(Msg, nil)))
	logger.Info("LogsController/RemoveAll", Msg)
}
