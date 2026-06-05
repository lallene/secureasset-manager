package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients map[*websocket.Conn]bool
	Mutex   sync.Mutex
}

var WS = &Hub{
	Clients: make(map[*websocket.Conn]bool),
}

func Broadcast(message string) {
	WS.Mutex.Lock()
	defer WS.Mutex.Unlock()

	for client := range WS.Clients {
		client.WriteJSON(map[string]string{
			"type":    "notification",
			"message": message,
		})
	}
}
