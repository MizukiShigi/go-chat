package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"mychat/internal/infrastructure/redis"
	"mychat/internal/infrastructure/websocket"
)

func serveWs(hub *websocket.Hub, presence *redis.UserPresence, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	fmt.Println(userID)
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &websocket.Client{
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		Presence: presence,
	}
	client.Hub.Register <- client

	ctx := context.Background()

	if err := presence.SetOnline(ctx, userID); err != nil {
		log.Println(err)
		return
	}

	go client.MaintainPresence(ctx)
	go client.WritePump()
	go client.ReadPump()
}

func main() {
	redisClient := redis.NewClient()
	presence := redis.NewUserPresence(redisClient)

	hub := websocket.NewHub(presence)
	go hub.Run()

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, presence, w, r)
	})

	log.Println("Server starting on :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
