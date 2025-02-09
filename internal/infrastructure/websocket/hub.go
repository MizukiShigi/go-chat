package websocket

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"mychat/internal/domain"
	"mychat/internal/infrastructure/redis"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Presence   *redis.UserPresence
}

func NewHub(presence *redis.UserPresence) *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Presence:   presence,
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
		case chat := <-h.Broadcast:
			h.broadcastChat(chat)
		case <-ticker.C:
			h.broadcastOnlineUsers()
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.Clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.Clients[client]; ok {
		delete(h.Clients, client)
		close(client.Send)
	}
}

// メッセージ送信
func (h *Hub) broadcast(message []byte) {
	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, client)
		}
	}
}

func (h *Hub) broadcastChat(chat []byte) {
	message := domain.Message{
		Type: domain.TypeChat,
		Content: domain.ChatConten{
			Chat: string(chat),
		},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal chat: %v", err)
		return
	}

	h.broadcast(messageBytes)
}

// オンラインユーザー一覧を送信
func (h *Hub) broadcastOnlineUsers() {
	users, err := h.Presence.GetOnlineUsers(context.Background())
	if err != nil {
		log.Printf("Failed to get online users: %v", err)
		return
	}

	message := domain.Message{
		Type: domain.TypePresence,
		Content: domain.PresenceContent{
			OnlineUsers: users,
		},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal online users: %v", err)
		return
	}

	h.broadcast(messageBytes)
}
