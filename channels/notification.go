package channels

import "github.com/gorilla/websocket"

type NotificationChannel struct {
	user string
	peer *websocket.Conn
}
