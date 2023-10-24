package session

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Store interface {
	CreateTokenSession(ctx context.Context, arg CreateTokenSessionParams) (*TokenSession, error)
	GetTokenSession(ctx context.Context, sessionId uuid.UUID) (*TokenSession, error)
	DeleteTokenSession(ctx context.Context, sessionId uuid.UUID) error
}

func NewSession(address string) (Store, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
	err := client.Ping(context.Background())
	return &Queries{
		db: *client,
	}, err.Err()
}
