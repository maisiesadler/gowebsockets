package wsmanager

import (
	"encoding/json"

	"github.com/rs/xid"
)

func (user *User) read() {
	for {
		select {
		case userMessage, ok := <-user.Send:
			if !ok {
				// todo: ?
			}

			recordedID := userToID[user]
			msgBytes, err := json.Marshal(userMessage)
			if err == nil {
				if socket, ok := connectedAuthIDToSocketID[recordedID]; ok {
					socket.send <- msgBytes
				}
			}
		}
	}
}

var connectedSocketIDToAuthID = make(map[string]string)
var connectedAuthIDToSocketID = make(map[string]*connectedSocket)

var loggedInIds = make(map[string]*User)
var userToID = make(map[*User]string)

func addKeyToConnectedSession(connectedSocket *connectedSocket, authID string) (isLoggedIn bool, user *User) {
	// todo: clean up?
	connectedSocketIDToAuthID[connectedSocket.ID] = authID
	connectedAuthIDToSocketID[authID] = connectedSocket
	user, isLoggedIn = loggedInIds[authID]
	return isLoggedIn, user
}

func createUser(authID string) *User {
	user := &User{
		ID:                xid.New(),
		Send:              make(chan *UserMessage),
		Receive:           make(chan *UserMessage),
		ReconnectedSocket: make(chan bool),
	}
	loggedInIds[authID] = user
	userToID[user] = authID
	return user
}

func authenticatedSocket(connectedSocketID string) (hasUser bool, user *User) {
	if authID, authed := connectedSocketIDToAuthID[connectedSocketID]; authed {
		user, hasUser := loggedInIds[authID]
		return hasUser, user
	}

	return false, nil
}

func socketDisconnected(socket *connectedSocket) {
	delete(connectedSocketIDToAuthID, socket.ID)
}
