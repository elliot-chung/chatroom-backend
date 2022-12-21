package main

import (
	"chatroom/internal"
	"log"
	"net/http"
)

const (
	CONN_HOST = ""
	CONN_PORT = "443"
)

func main() {
	bh := internal.NewBroadcastHub()
	canvas := internal.NewCanvas()
	go internal.CanvasCleaner(canvas)
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
