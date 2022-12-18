package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var originWhitelist = []string{"http://localhost:5173",
	"http://localhost:3334", "http://localhost:3335"}

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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		defer conn.Close()

		for {
			// Read message from browser
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			fmt.Println(string(msg))
		}
	})

	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT + "...")
	http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)
}
