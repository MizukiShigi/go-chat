package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	keySessionPrefix = "session:"
	SessionTTL       = 24 * time.Hour
)

type Session struct {
	Client *Client
}

func NewSession(client *Client) *Session {
	return &Session{Client: client}
}

func (s *Session) CreateSession(ctx context.Context, userID string) (string, error) {
	sessionID := uuid.New().String()
	key := keySessionPrefix + sessionID

	err := s.Client.Set(ctx, key, userID, SessionTTL).Err()
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *Session) GetSession(ctx context.Context, sessionID string) (string, error) {
	key := keySessionPrefix + sessionID
	userID, err := s.Client.Get(ctx, key).Result()
	if err != nil {
        return "", err
    }

	s.Client.Expire(ctx, key, SessionTTL)

	return userID, nil
}
