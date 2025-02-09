package redis

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const keyTypingPrefix = "typing:"

type TypingNotification struct {
	Client *Client
}

func NewTypingNotification(client *Client) *TypingNotification {
	return &TypingNotification{Client: client}
}

func (tn *TypingNotification) SetTyping(ctx context.Context, roomID, userID string) error {
	key := fmt.Sprintf(keyTypingPrefix+"%s:%s", roomID, userID)
	return tn.Client.Set(ctx, key, true, 3*time.Second).Err()
}

func (tn *TypingNotification) GetTypingUsers(ctx context.Context, roomID string) ([]string, error) {
	pattern := fmt.Sprintf(keyTypingPrefix+"%s:*", roomID)
	keys, err := tn.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	users := make([]string, len(keys))
	for i, key := range keys {
		// "typing:roomID:userID"
		parts := strings.Split(key, ":")
		if len(parts) == 3 {
			users[i] = parts[2]
		}
	}

	return users, nil
}
