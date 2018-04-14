package wsmanager

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

type ConnectedSocket struct {
	ID          string
	socket      *websocket.Conn
	send        chan []byte
	UserCreated chan *User
}

type User struct {
	ID                xid.ID
	Name              string
	Send              chan *UserMessage
	Receive           chan *UserMessage
	ReconnectedSocket chan bool
}

type UserMessage struct {
	Command string
	Params  map[string]string
}

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
				if socket, ok := connectedSessionIDToSocketID[recordedID]; ok {
					socket.send <- msgBytes
				}
			}
		}
	}
}

var connectedSocketIDToSessionID = make(map[string]string)
var connectedSessionIDToSocketID = make(map[string]*ConnectedSocket)

var loggedInIds = make(map[string]*User)
var userToID = make(map[*User]string)

func addKeyToConnectedSession(connectedSocket *ConnectedSocket, recordedID string) (isLoggedIn bool, user *User) {
	// todo: clean up?
	connectedSocketIDToSessionID[connectedSocket.ID] = recordedID
	connectedSessionIDToSocketID[recordedID] = connectedSocket
	user, isLoggedIn = loggedInIds[recordedID]
	if isLoggedIn {
		user.ReconnectedSocket <- true
	}
	return isLoggedIn, user
}

func logIn(recordedID string, username string) *User {
	user := &User{
		Name:              username,
		ID:                xid.New(),
		Receive:           make(chan *UserMessage),
		Send:              make(chan *UserMessage),
		ReconnectedSocket: make(chan bool),
	}
	loggedInIds[recordedID] = user
	userToID[user] = recordedID
	go user.read()
	return user
}

func authenticatedSocket(connectedSocketID string) (authed bool, recordedID string) {
	recordedID, ok := connectedSocketIDToSessionID[connectedSocketID]
	return ok, recordedID
}

func sessionIDIsLoggedIn(recordedID string) (loggedIn bool, user *User) {
	if user, ok := loggedInIds[recordedID]; ok {
		return true, user
	}
	return false, &User{}
}

func socketDisconnected(socket *ConnectedSocket) {
	delete(connectedSessionIDToSocketID, socket.ID)
}
