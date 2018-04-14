package wsmanager

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

func Create(conn *websocket.Conn) *ConnectedSocket {
	client := &ConnectedSocket{
		ID:          xid.New().String(),
		socket:      conn,
		send:        make(chan []byte),
		UserCreated: make(chan *User),
	}
	go client.read()
	go client.write()
	msg := getSendAuthUserMessage()
	msgBytes, _ := json.Marshal(msg)
	client.send <- msgBytes

	return client
}

func (c *ConnectedSocket) read() {
	defer func() {
		c.disconnected()
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			break
		}
		userMessage := handleMsg(c, string(message))
		if userMessage != nil {
			msgBytes, _ := json.Marshal(userMessage)
			c.send <- msgBytes
		}
	}
}

func (c *ConnectedSocket) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *ConnectedSocket) disconnected() {
	// logger.Log("wsmanagerwsmanager", "disconnected", "socket disconnected")
	socketDisconnected(c)
}
