package server

import (
	"encoding/json"
)

func handleMsg(connectedSocket *connectedSocket, strMsg string) {
	userMessage := &UserMessage{}
	err := json.Unmarshal([]byte(strMsg), userMessage)
	if err != nil {
		// logger.Log("session", "handleMsg", "unmarshal")
		panic(err)
	}
	if userMessage.Command == "auth" {
		authID := userMessage.Params["authID"]
		if isLoggedIn, user := addKeyToConnectedSession(connectedSocket, authID); isLoggedIn {
			user.ReconnectedSocket <- true
		} else {
			newUser := createUser(authID)
			connectedSocket.UserCreated <- newUser
		}
	} else {
		userExists, user := authenticatedSocket(connectedSocket.ID)

		if userExists {
			user.Receive <- userMessage
		} else {
			go sendAuthUserMessage(connectedSocket)
		}
	}
}
