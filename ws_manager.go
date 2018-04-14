package wsmanager

//https://www.thepolyglotdeveloper.com/2016/12/create-real-time-chat-app-golang-angular-2-websockets/

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var newUsers = make(chan *User)

// Start starts the websocket server and returns a pointer to a channel returning new users
func Start() chan *User {
	http.HandleFunc("/echo", wsPage)
	http.ListenAndServe(":8080", nil)

	return newUsers
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	client := Create(conn)
	go waitForUser(client)
}

func waitForUser(socket *ConnectedSocket) {
	user := <-socket.UserCreated
	newUsers <- user
}
