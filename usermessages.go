package server

import (
	"encoding/json"

	"github.com/rs/xid"
)

func sendAuthUserMessage(connectedSocket *connectedSocket) {
	userMessage := UserMessage{
		Command: "sendAuth",
		Params:  map[string]string{"id": xid.New().String()},
	}

	msgBytes, _ := json.Marshal(userMessage)
	connectedSocket.send <- msgBytes
}
