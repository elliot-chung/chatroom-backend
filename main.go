package main

import (
	"chatroom/internal"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
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

	http.HandleFunc("/chatroom", func(w http.ResponseWriter, r *http.Request) {
		internal.RequestHandler(canvas, bh, w, r)
	})

	// creates a new file watcher for App_offline.htm
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("ERROR", err)
	}
	defer watcher.Close()

	// watch for App_offline.htm and exit the program if present
	// This allows continuous deployment on App Service as the .exe will not be
	// terminated otherwise
	go func() {
		for event := range watcher.Events {
			if strings.HasSuffix(event.Name, "app_offline.htm") {
				log.Println("Exiting due to app_offline.htm being present")
				os.Exit(0)
			}
		}
	}()

	// get the current working directory and watch it
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("ERROR", err)
	}

	if err := watcher.Add(currentDir); err != nil {
		log.Println("ERROR", err)
	}

	// Azure env variables
	hostname := os.Getenv("WEBSITE_HOSTNAME")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on " + hostname + ":" + port)
	err = http.ListenAndServe(hostname+":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
