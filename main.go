package main

import (
	"chatroom/internal"
	"fmt"
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
		internal.RequestHandler(canvas, bh, w, r)
	})

	// creates a new file watcher for App_offline.htm
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	// watch for App_offline.htm and exit the program if present
	// This allows continuous deployment on App Service as the .exe will not be
	// terminated otherwise
	go func() {
		for event := range watcher.Events {
			if strings.HasSuffix(event.Name, "app_offline.htm") {
				fmt.Println("Exiting due to app_offline.htm being present")
				os.Exit(0)
			}
		}
	}()

	// get the current working directory and watch it
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("ERROR", err)
	}

	if err := watcher.Add(currentDir); err != nil {
		fmt.Println("ERROR", err)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on " + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
