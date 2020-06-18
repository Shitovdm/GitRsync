package controllers

import (
	"fmt"
	"github.com/Shitovdm/GitRsync/src/components/helpers"
	"github.com/Shitovdm/GitRsync/src/components/interface"
	"github.com/Shitovdm/GitRsync/src/components/logger"
	"github.com/Shitovdm/GitRsync/src/models"
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
	var subsctibeToLogRequest models.RuntimeLogsRequest
	conn, err := helpers.WsHandler(c.Writer, c.Request, &subsctibeToLogRequest)
	if err != nil {
		logger.Error("LogsController/Subscribe", err.Error())
		return
	}

	switch subsctibeToLogRequest.Action {
	case "init":
		_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\x1b[92m%s\x1b[0m", " GitRsync (C) Shitov Dmitry")))
		for _, logNote := range logger.GetRuntimeLogs() {
			//	Only fot current session.
			if logNote.SessionID == logger.GetSessionId() {
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

func (ctrl LogsController) RemoveRuntime(c *gin.Context) {

	var removeRuntimeLogsRequest models.RuntimeLogsRequest
	conn, err := helpers.WsHandler(c.Writer, c.Request, &removeRuntimeLogsRequest)
	if err != nil {
		Msg := fmt.Sprintf("Error while removing runtime logs! %s", err.Error())
		logger.Error("LogsController/RemoveRuntime", Msg)
		return
	}

	err = logger.ClearRuntimeLogs()
	if err != nil {
		Msg := fmt.Sprintf("%s", err.Error())
		logger.Error("LogsController/RemoveRuntime", Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(Msg)))
		return
	}

	Msg := "Runtime logs successfully removed!"
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	logger.Info("LogsController/RemoveRuntime", Msg)
	return
}

func (ctrl LogsController) RemoveAll(c *gin.Context) {

	var removeAllLogsRequest models.RuntimeLogsRequest
	conn, err := helpers.WsHandler(c.Writer, c.Request, &removeAllLogsRequest)
	if err != nil {
		Msg := fmt.Sprintf("Error while removing all logs! %s", err.Error())
		logger.Error("LogsController/RemoveAll", Msg)
		return
	}

	err = logger.ClearAllLogs()
	if err != nil {
		Msg := fmt.Sprintf("%s", err.Error())
		logger.Error("LogsController/RemoveAll", Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonError(Msg)))
		return
	}

	Msg := "All logs successfully removed!"
	_ = conn.WriteMessage(websocket.TextMessage, []byte(BuildWsJsonSuccess(Msg, nil)))
	logger.Info("LogsController/RemoveAll", Msg)
	return
}
