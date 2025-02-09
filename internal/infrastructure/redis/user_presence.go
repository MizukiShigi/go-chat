package redis

import (
	"context"
	"fmt"
	"time"
)

const (
	keyOnlinePrefix    = "online:"
	PresenseTTLSeconds = 30
)

type UserPresence struct {
	Client *Client
}

func NewUserPresence(client *Client) *UserPresence {
	return &UserPresence{Client: client}
}

func (up *UserPresence) SetOnline(ctx context.Context, userID string) error {
	key := fmt.Sprintf(keyOnlinePrefix+"%s", userID)
	return up.Client.Set(ctx, key, true, PresenseTTLSeconds*time.Second).Err()
}

func (up *UserPresence) UpdatePresence(ctx context.Context, userID string) error {
	return up.SetOnline(ctx, userID)
}

func (up *UserPresence) SetOffline(ctx context.Context, userID string) error {
	key := fmt.Sprintf(keyOnlinePrefix+"%s", userID)
    return up.Client.Del(ctx, key).Err()
}

func (up *UserPresence) IsOnline(ctx context.Context, userID string) (bool, error) {
	key := fmt.Sprintf(keyOnlinePrefix+"%s", userID)
	result, err := up.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (up *UserPresence) GetOnlineUsers(ctx context.Context) ([]string, error) {
	keys, err := up.Client.Keys(ctx, keyOnlinePrefix+"*").Result()
	if err != nil {
		return nil, err
	}
	users := make([]string, 0, len(keys))
	for _, key := range keys {
		users = append(users, key[len(keyOnlinePrefix):])
	}
	return users, nil
}
