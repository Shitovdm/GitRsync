package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// MarshalFunc returns byte data.
type MarshalFunc func(v interface{}) ([]byte, error)

// UnmarshalFunc gets byte data.
type UnmarshalFunc func(data []byte, v interface{}) error

// format returns new UUID in v4 format.
type format struct {
	m  MarshalFunc
	um UnmarshalFunc
}

// WebSocketUpgrade buffer size.
var WebSocketUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WsHandler handles stream messages.
func WsHandler(w http.ResponseWriter, r *http.Request, body interface{}) (*websocket.Conn, error) {
	conn, err := WebSocketUpgrade.Upgrade(w, r, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to set websocket upgrade! %s", err.Error()))
	}

	_, data, _ := conn.ReadMessage()
	cFormat := format{m: json.Marshal, um: json.Unmarshal}
	if err := cFormat.um(data, body); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to unmarshal form data! %s", err.Error()))
	}

	return conn, nil
}
