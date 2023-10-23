package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	db "github.com/punkzberryz/todo/db/sqlc"
	"github.com/punkzberryz/todo/token"
)

// Create new AccessToken & RefreshToken
// used for CreateNewUser or UserLogin
type newTokenResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (*newTokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *Server) createNewAccessToken(user *db.User, r *http.Request) (*newTokenResponse, error) {
	userPayload := token.User{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(userPayload, server.config.AccessTokenDuration)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(userPayload, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	session, err := server.store.CreateSession(r.Context(), db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       refreshPayload.User.ID,
		RefreshToken: refreshToken,
		UserAgent:    r.UserAgent(),
		ClientIp:     r.RemoteAddr,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, err
	}
	rsp := newTokenResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}
	return &rsp, nil
}

// For refresh token
type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (c *renewAccessTokenRequest) Bind(r *http.Request) error {
	if c.RefreshToken == "" {
		return fmt.Errorf("missing refresh token")
	}
	return nil
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (*renewAccessTokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *Server) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	data := &renewAccessTokenRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(data.RefreshToken)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	session, err := server.store.GetSession(r.Context(), refreshPayload.ID)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	if session.UserID != refreshPayload.User.ID {
		err := fmt.Errorf("incorrect session user")
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	if session.RefreshToken != data.RefreshToken {
		err := fmt.Errorf("mismatch session token")
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	userPayload := token.User{
		ID:       refreshPayload.User.ID,
		Email:    refreshPayload.User.Email,
		Username: refreshPayload.User.Username,
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(userPayload, server.config.AccessTokenDuration)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	if err := render.Render(w, r, &rsp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// Delete token sesion
type removeTokenSessionRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (c *removeTokenSessionRequest) Bind(r *http.Request) error {
	if c.RefreshToken == "" {
		return fmt.Errorf("missing refresh token")
	}
	return nil
}

type removeTokenSessionResponse struct {
	Message string `json:"message"`
}

func (*removeTokenSessionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Caveat is that we only delete Refresh Token Session,
// This means user cannot refresh AccessToken when the token is expired
// However, user can still
// perform some tasks using the AccessToken, if the token is not expired
func (server *Server) removeTokenSession(w http.ResponseWriter, r *http.Request) {
	data := &removeTokenSessionRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(data.RefreshToken)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	//Delete session
	err = server.store.DeleteSession(r.Context(), refreshPayload.ID)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	rsp := &removeTokenSessionResponse{
		Message: fmt.Sprintf("%s logout successfully", refreshPayload.User.Email),
	}
	if err := render.Render(w, r, rsp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}
