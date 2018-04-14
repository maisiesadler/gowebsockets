# gowebsocket
Web socket manager written in go

When a new socket connects it must send an auth id, which will connect it to a wsmanager.User, which will then be sent to the newUsers channel.
When the socket disconnects and reconnects it must resend the same auth id to be reconnected to the same wsmanager.User.

Example

```go
func main() {
	newUsers := wsmanager.Start()
	go waitForUsers(newUsers)
}

func waitForUsers(users chan *wsmanager.User) {
	user := <-users
	go handleUser(user)
}

func handleUser(user *wsmanager.User) {
	for {
		select {
		case msg, ok := <-user.Receive:
			if !ok {
				break
			}
			response := handleAuthenticatedMsg(user, msg)
			user.Send <- response
		case isLoggedIn, ok := <-user.ReconnectedSocket:
			if !ok {
				break
			}
			if isLoggedIn {
				response := handleAuthenticatedMsg(user, &wsmanager.UserMessage{Command: "tryrejoin"})
				user.Send <- response
			}
		}
	}
}
```
