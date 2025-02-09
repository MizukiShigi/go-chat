package websocket

import (
	"context"
	"log"
	"time"

	"mychat/internal/infrastructure/redis"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	Presence *redis.UserPresence
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Presence.SetOffline(context.Background(), c.UserID)
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.Hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)
		if err := w.Close(); err != nil {
			return
		}
	}
}

func (c *Client) MaintainPresence(ctx context.Context) {
	// redisのttlより少し短めに設定
	ticker := time.NewTicker((redis.PresenseTTLSeconds - 5) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println(c.UserID)
			if err := c.Presence.UpdatePresence(ctx, c.UserID); err != nil {
				log.Printf("failed to update presence: %v", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
