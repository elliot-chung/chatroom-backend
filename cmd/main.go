package main

import (
	"chatroom/internal"
	"log"
	"net/http"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	bh := internal.NewBroadcastHub()
	canvas := internal.NewCanvas()
	go bh.Run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		internal.RequestHandler(canvas, bh, w, r)
	})

	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT + "...")
	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
