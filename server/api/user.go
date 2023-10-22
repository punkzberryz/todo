package api

import (
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	db "github.com/punkzberryz/todo/db/sqlc"
	"github.com/punkzberryz/todo/token"
	"github.com/punkzberryz/todo/util"
)

// For create new user request
type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (c *createUserRequest) Bind(r *http.Request) error {
	if c.Email == "" || c.Username == "" || c.Password == "" {
		return fmt.Errorf("missing email, username or/add password fields")
	}
	if len(c.Password) < 6 {
		return fmt.Errorf("password must be at least 6 letters")
	}
	_, err := mail.ParseAddress(c.Email)
	if err != nil {
		return fmt.Errorf("email is invalid")
	}
	return nil
}

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (*userResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type loginOrRegisterUserResponse struct {
	SessionID             uuid.UUID     `json:"session_id"`
	AccessToken           string        `json:"access_token"`
	AccessTokenExpiresAt  time.Time     `json:"access_token_expires_at"`
	RefreshToken          string        `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time     `json:"refresh_token_expires_at"`
	User                  *userResponse `json:"user"`
}

func (*loginOrRegisterUserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *Server) createUser(w http.ResponseWriter, r *http.Request) {
	data := &createUserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       data.Username,
		HashedPassword: hashedPassword,
		Email:          data.Email,
	}

	user, err := server.store.CreateUser(r.Context(), arg)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	userPayload := token.User{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	token, payload, err := server.tokenMaker.CreateToken(userPayload, server.config.AccessTokenDuration)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	userRsp := userResponse{
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	rsp := loginOrRegisterUserResponse{
		AccessToken:          token,
		SessionID:            payload.ID,
		AccessTokenExpiresAt: payload.ExpiredAt,
		User:                 &userRsp,
	}
	if err := render.Render(w, r, &rsp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// For login user request
type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *loginUserRequest) Bind(r *http.Request) error {
	if c.Email == "" || c.Password == "" {
		return fmt.Errorf("missing email or/add password fields")
	}
	if len(c.Password) < 6 {
		return fmt.Errorf("password must be at least 6 letters")
	}
	_, err := mail.ParseAddress(c.Email)
	if err != nil {
		return fmt.Errorf("email is invalid")
	}
	return nil
}

func (server *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	data := &loginUserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	arg := db.GetUserParams{
		Email: data.Email,
	}
	user, err := server.store.GetUser(r.Context(), arg)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("email or password is incorrect")))
		return
	}

	err = util.ComparePassword(data.Password, user.HashedPassword)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(fmt.Errorf("email or password is incorrect")))
		return
	}
	userPayload := token.User{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	token, accessPayload, err := server.tokenMaker.CreateToken(userPayload, server.config.AccessTokenDuration)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(userPayload, server.config.RefreshTokenDuration)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
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
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	userRsp := userResponse{
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	rsp := loginOrRegisterUserResponse{
		SessionID:             session.ID,
		AccessToken:           token,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  &userRsp,
	}
	if err := render.Render(w, r, &rsp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// get current user from JWT
func (server *Server) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(payloadKey).(*token.Payload)
	arg := db.GetUserParams{
		ID: payload.User.ID,
	}
	user, err := server.store.GetUser(r.Context(), arg)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	rsp := userResponse{
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	if err := render.Render(w, r, &rsp); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}
