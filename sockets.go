package server

import (
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

func create(conn *websocket.Conn) *connectedSocket {
	client := &connectedSocket{
		ID:          xid.New().String(),
		socket:      conn,
		send:        make(chan []byte),
		UserCreated: make(chan *User),
	}
	go client.read()
	go client.write()
	sendAuthUserMessage(client)

	return client
}

func (c *connectedSocket) read() {
	defer func() {
		c.disconnected()
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			break
		}
		handleMsg(c, string(message))
	}
}

func (c *connectedSocket) write() {
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

func (c *connectedSocket) disconnected() {
	// logger.Log("session", "disconnected", "socket disconnected")
	socketDisconnected(c)
}
