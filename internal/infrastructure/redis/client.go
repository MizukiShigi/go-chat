package redis

import "github.com/redis/go-redis/v9"

type Client struct {
	*redis.Client
}

func NewClient() *Client {
	return &Client{
		redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}
