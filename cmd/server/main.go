package main

import (
	"mychat/internal/infrastructure/websocket"
	"log"
	"net/http"
)

func serveWs(hub *websocket.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &websocket.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func main() {
    hub := websocket.NewHub()
    go hub.Run()

    fs := http.FileServer(http.Dir("web/static"))
    http.Handle("/", fs)

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })

    log.Println("Server starting on :8081")
    err := http.ListenAndServe(":8081", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
