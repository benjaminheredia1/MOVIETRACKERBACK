package redis

import (
	"github.com/redis/go-redis/v9"
)

func NewConnection(host, port, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
	})

	return client
}
