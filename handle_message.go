package wsmanager

import (
	"encoding/json"

	"github.com/rs/xid"
)

type messageHandler func(connectedSocket *ConnectedSocket, params map[string]string) (reply *UserMessage)

var handlers = make(map[string]messageHandler)

func registerHandler(match string, handler messageHandler) {
	if _, ok := handlers[match]; ok {
		panic("Handler already registered for " + match)
	}

	handlers[match] = handler
}

func getSendAuthUserMessage() *UserMessage {
	return &UserMessage{
		Command: "sendAuth",
		Params:  map[string]string{"id": xid.New().String()},
	}
}

func getLoggedInUserMessage(user *User) *UserMessage {
	return &UserMessage{
		Command: "loggedin",
		Params:  map[string]string{"username": user.Name},
	}
}

func init() {

	registerHandler("login", func(connectedSocket *ConnectedSocket, params map[string]string) *UserMessage {
		authedSocket, recordedID := authenticatedSocket(connectedSocket.ID)
		if !authedSocket {
			return getSendAuthUserMessage()
		}
		username := params["username"]
		user := logIn(recordedID, username)
		connectedSocket.UserCreated <- user
		return getLoggedInUserMessage(user)
	})
	registerHandler("auth", func(connectedSocket *ConnectedSocket, params map[string]string) *UserMessage {
		sessionid := params["sessionid"]
		if isLoggedIn, user := addKeyToConnectedSession(connectedSocket, sessionid); isLoggedIn {
			return getLoggedInUserMessage(user)
		}

		return nil
	})
}

func handleMsg(connectedSocket *ConnectedSocket, strMsg string) *UserMessage {
	userMessage := &UserMessage{}
	err := json.Unmarshal([]byte(strMsg), userMessage)
	if err != nil {
		panic(err)
		// logger.Log("session", "handleMsg", "unmarshal")
	}
	if handler, match := handlers[userMessage.Command]; match {
		reply := handler(connectedSocket, userMessage.Params)
		return reply
	}

	authedSocket, recordedID := authenticatedSocket(connectedSocket.ID)

	if authedSocket {
		isLoggedIn, user := sessionIDIsLoggedIn(recordedID)
		if isLoggedIn {
			user.Receive <- userMessage
			return nil
		}
		return &UserMessage{Command: "plslogin"}
	}

	return getSendAuthUserMessage()
}
