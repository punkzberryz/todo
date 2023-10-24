package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrTokenSessionNotFound = fmt.Errorf("token session not found")
)

type Queries struct {
	db redis.Client
}

type CreateTokenSessionParams struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"userId"`
	RefreshToken string    `json:"refreshToken"`
	UserAgent    string    `json:"userAgent"`
	ClientIp     string    `json:"clientIp"`
	IsBlocked    bool      `json:"isBlocked"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type TokenSession struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"userId"`
	RefreshToken string    `json:"refreshToken"`
	UserAgent    string    `json:"userAgent"`
	ClientIp     string    `json:"clientIp"`
	IsBlocked    bool      `json:"isBlocked"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (q *Queries) CreateTokenSession(ctx context.Context, arg CreateTokenSessionParams) (*TokenSession, error) {
	id := arg.ID.String()

	token := TokenSession{
		ID:           arg.ID,
		UserID:       arg.UserID,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIp,
		IsBlocked:    arg.IsBlocked,
		ExpiresAt:    arg.ExpiresAt,
		CreatedAt:    time.Now(),
	}

	json, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	err = q.db.Set(ctx, id, json, time.Until(arg.ExpiresAt)).Err()
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (q *Queries) GetTokenSession(ctx context.Context, sessionId uuid.UUID) (*TokenSession, error) {
	id := sessionId.String()
	var token = &TokenSession{}
	result, err := q.db.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrTokenSessionNotFound
		}
		return nil, err
	}
	if err := json.Unmarshal([]byte(result), token); err != nil {
		return nil, err
	}
	return token, nil
}

func (q *Queries) DeleteTokenSession(ctx context.Context, sessionId uuid.UUID) error {
	id := sessionId.String()
	err := q.db.Del(ctx, id).Err()
	if err != nil {
		return err
	}
	return nil
}
