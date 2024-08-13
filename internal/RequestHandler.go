package internal

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Request struct {
	Type  string `json:"type"`
	User  string `json:"user"`
	Text  string `json:"text,omitempty"`
	Color int    `json:"color"`
	X     int    `json:"x,omitempty"`
	Y     int    `json:"y,omitempty"`
}

var originWhitelist = []string{
	"http://localhost:5174",
	"https://chatroom-frontend-one.vercel.app",
	"https://chatroom-frontend-ec-2018.vercel.app",
	"https://chatroom-frontend-git-main-ec-2018.vercel.app",
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		for _, origin := range originWhitelist {
			if r.Header.Get("Origin") == origin {
				return true
			}
		}
		return false
	},
}

func RequestHandler(canvas *Canvas, hub *BroadcastHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer conn.Close()

	// Register our new client
	hub.register <- conn
	log.Println("connection established:", conn.RemoteAddr())
	conn.SetCloseHandler(func(code int, text string) error {
		log.Println("connection closed:", conn.RemoteAddr())
		hub.unregister <- conn
		return nil
	})

	// Send canvas to client
	data, err := canvas.MarshalJSON()
	if err != nil {
		log.Println(err)
		hub.unregister <- conn
		return
	}
	conn.WriteMessage(websocket.TextMessage, data)

	for {
		// Read message as json
		var msg Request
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure) {
				break
			}
			log.Println("Read Error:", err)
			break
		}

		// Handle different message types
		switch msg.Type {
		case "message":
			log.Println("Message from", conn.RemoteAddr(), "("+msg.User+")", ":", msg.Text)
			// Broadcast message to all clients
			jsonMessage, err := json.Marshal(struct {
				Type  string
				User  string
				Text  string
				Color int
			}{
				Type:  msg.Type,
				User:  msg.User,
				Text:  msg.Text,
				Color: msg.Color,
			})

			if err != nil {
				log.Println(err)
				hub.unregister <- conn
				break
			}

			hub.broadcast <- []byte(jsonMessage)

		case "draw":
			// Draw on canvas
			log.Println(conn.RemoteAddr(), "drawing at", msg.X, msg.Y, "with color", msg.Color)
			canvas.SetCoordinate(msg.X, msg.Y, msg.Color)
			data, err := canvas.MarshalJSON()

			if err != nil {
				log.Println(err)
				hub.unregister <- conn
				break
			}

			// Broadcast canvas to all clients
			hub.broadcast <- data
		}
	}
}
