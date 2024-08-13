package main

import (
	"chatroom/internal"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	bh := internal.NewBroadcastHub()
	canvas := internal.NewCanvas()
	go internal.CanvasCleaner(canvas)
	go bh.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// hello world, the web server
		w.Write([]byte("Hello, World!"))
	})

	http.HandleFunc("/connections", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%d", bh.ConnectionCount())))
		log.Println(bh.ConnectionCount())
	})

	http.HandleFunc("/chatroom", func(w http.ResponseWriter, r *http.Request) {
		internal.RequestHandler(canvas, bh, w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
