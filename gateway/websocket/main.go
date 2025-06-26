package main

import (
	"log"
	"net/http"
)

func main() {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		ServeWs(hub, writer, request)
	})

	log.Println("WebSocket Gateway started on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
