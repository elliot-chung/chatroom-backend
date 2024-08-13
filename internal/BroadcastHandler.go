package internal

import (
	"time"

	"github.com/gorilla/websocket"
)

type BroadcastHub struct {
	// Registered connections.
	connections map[*websocket.Conn]bool

	// Outbound message to the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *websocket.Conn

	// Unregister requests from connections.
	unregister chan *websocket.Conn
}

func NewBroadcastHub() *BroadcastHub {
	return &BroadcastHub{
		broadcast:   make(chan []byte),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
		connections: make(map[*websocket.Conn]bool),
	}
}

func (h *BroadcastHub) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
		case m := <-h.broadcast:
			time.Sleep(10 * time.Millisecond) // 10 ms rebroadcast delay
			for c := range h.connections {
				c.WriteMessage(websocket.TextMessage, m)
			}
		}
	}
}
