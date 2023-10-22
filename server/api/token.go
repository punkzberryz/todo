package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/punkzberryz/todo/token"
)

// For login user request
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
