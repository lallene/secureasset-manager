package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)

	if err != nil {
		return
	}

	WS.Mutex.Lock()
	WS.Clients[conn] = true
	WS.Mutex.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			WS.Mutex.Lock()
			delete(WS.Clients, conn)
			WS.Mutex.Unlock()
			conn.Close()
			break
		}
	}
}
