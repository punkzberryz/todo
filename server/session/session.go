package session

import (
	"github.com/go-redis/redis/v8"
)

func NewSession(address string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
	return client
}
