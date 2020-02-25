package Helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type MarshalFunc func(v interface{}) ([]byte, error)

type UnmarshalFunc func(data []byte, v interface{}) error

type format struct {
	m  MarshalFunc
	um UnmarshalFunc
}

var WebSocketUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsHandler(w http.ResponseWriter, r *http.Request, body interface{}) error {
	conn, err := WebSocketUpgrade.Upgrade(w, r, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to set websocket upgrade! %s", err.Error()))
	}

	_, data, _ := conn.ReadMessage()
	cFormat := format{m: json.Marshal, um: json.Unmarshal}
	if err := cFormat.um(data, body); err != nil {
		return errors.New(fmt.Sprintf("Failed to unmarshal form data! %s", err.Error()))
	}

	return nil
}
