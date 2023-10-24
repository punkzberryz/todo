package token

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	db "github.com/punkzberryz/todo/db/sqlc"
)

var (
	ErrBlockedSession       = fmt.Errorf("blocked session")
	ErrIncorrectSessionUser = fmt.Errorf("incorrect session user")
	ErrMismatchSessionToken = fmt.Errorf("mismatch session token")
	ErrSessionExpired       = fmt.Errorf("expired session")
)

type Token struct {
	Store db.Store
	Maker
	RefreshTokenDuration time.Duration
	AccessTokenDuration  time.Duration
}

// Create new AccessToken & RefreshToken
// used for CreateNewUser or UserLogin
type NewTokenResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type CreateTokenParams struct {
	User
	UserAgent string
	ClientIp  string
}

func (t Token) CreateNewAccessToken(ctx context.Context, arg CreateTokenParams) (*NewTokenResponse, error) {
	accessToken, accessPayload, err := t.Maker.CreateToken(arg.User, t.AccessTokenDuration)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshPayload, err := t.Maker.CreateToken(arg.User, t.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	session, err := t.Store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       refreshPayload.User.ID,
		RefreshToken: refreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, err
	}
	rsp := NewTokenResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}
	return &rsp, nil
}

// For renew the AccessToken with refreshToken
type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (t Token) RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenResponse, error) {
	refreshPayload, err := t.Maker.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}
	session, err := t.Store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		return nil, err
	}
	if session.IsBlocked {
		return nil, ErrBlockedSession
	}
	if session.UserID != refreshPayload.User.ID {
		return nil, ErrIncorrectSessionUser
	}

	if session.RefreshToken != refreshToken {
		return nil, ErrMismatchSessionToken
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	accessToken, accessPayload, err := t.Maker.CreateToken(refreshPayload.User, t.AccessTokenDuration)
	if err != nil {
		return nil, err
	}
	rsp := RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	return &rsp, nil
}

// Delete token session
func (t Token) DeleteTokenSession(ctx context.Context, refreshToken string) error {
	refreshPayload, err := t.Maker.VerifyToken(refreshToken)
	if err != nil {
		return err
	}
	err = t.Store.DeleteSession(ctx, refreshPayload.ID)
	return err
}
