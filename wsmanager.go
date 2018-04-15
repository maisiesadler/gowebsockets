package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Start starts the websocket server and returns a pointer to a channel returning new users
func Start(newUsers chan<- *User) {
	http.HandleFunc("/echo", func(res http.ResponseWriter, req *http.Request) {
		conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
		if error != nil {
			http.NotFound(res, req)
			return
		}
		client := create(conn)
		user := <-client.UserCreated
		newUsers <- user
	})
	http.ListenAndServe(":8080", nil)
}
