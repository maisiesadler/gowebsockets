package wsmanager

import (
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

type connectedSocket struct {
	ID          string
	socket      *websocket.Conn
	send        chan []byte
	UserCreated chan *User
}

type User struct {
	ID                xid.ID
	Send              chan *UserMessage
	Receive           chan *UserMessage
	ReconnectedSocket chan bool
}

type UserMessage struct {
	Command string
	Params  map[string]string
}
