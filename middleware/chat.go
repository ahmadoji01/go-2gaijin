package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/pkg/websocket"
)

func ServeWs(c *gin.Context) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(c.Writer, c.Request)
	if err != nil {
		fmt.Fprintf(c.Writer, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: Pool,
	}

	Pool.Register <- client
	client.Read()
}
