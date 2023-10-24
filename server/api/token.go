package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

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

	newToken, err := server.token.RenewAccessToken(r.Context(), data.RefreshToken)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          newToken.AccessToken,
		AccessTokenExpiresAt: newToken.AccessTokenExpiresAt,
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

	//Delete session
	err := server.token.DeleteTokenSession(r.Context(), data.RefreshToken)
	if err != nil {
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	rsp := &removeTokenSessionResponse{
		Message: "logout successfully",
	}
	if err := render.Render(w, r, rsp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}
