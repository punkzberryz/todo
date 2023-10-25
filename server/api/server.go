package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	db "github.com/punkzberryz/todo/db/sqlc"
	"github.com/punkzberryz/todo/service/auth"
	"github.com/punkzberryz/todo/service/mail"
	"github.com/punkzberryz/todo/service/task"
	"github.com/punkzberryz/todo/service/token"
	"github.com/punkzberryz/todo/session"
	"github.com/punkzberryz/todo/util"
)

type Server struct {
	config util.Config
	Router *chi.Mux
	auth   auth.Auth
	task   task.Task
	token  token.Token
	mail   mail.EmailSender
}

// Create new HTTP server and setup routing
func NewServer(config util.Config, store *db.Store, session *session.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	token := token.Token{
		Session:              *session,
		Maker:                tokenMaker,
		RefreshTokenDuration: config.RefreshTokenDuration,
		AccessTokenDuration:  config.AccessTokenDuration,
	}
	auth := auth.Auth{
		Store: *store,
	}
	task := task.Task{
		Store: *store,
	}
	mailSender := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	server := &Server{
		config: config,
		auth:   auth,
		task:   task,
		token:  token,
		mail:   mailSender,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome v2"))
	})

	//user-route
	r.Route("/user", func(r chi.Router) {
		r.Post("/", server.createUser)                                 //POST /user/
		r.Post("/login", server.loginUser)                             //POST /user/login
		r.Post("/logout", server.removeTokenSession)                   //POST /user/logout
		r.Post("/reset-password-request", server.resetPasswordRequest) //POST /user/reset-password-request
		r.Post("/reset-password", server.resetPassword)                //POST /user/reset-password
	})

	//token
	r.Post("/tokens/renew_access", server.renewAccessToken)

	// user-route-protected
	r.Route("/me", func(r chi.Router) {
		r.Use(server.authMiddleware)
		r.Get("/", server.getCurrentUser) //GET /me/
	})
	//task-route
	r.Route("/task", func(r chi.Router) {
		r.Use(server.authMiddleware)             //require Header {Authorization: Bearer token}
		r.Get("/{taskID}", server.getTask)       //GET /task/123
		r.Post("/", server.createTask)           //POST /task/123
		r.Get("/", server.getTaskList)           //GET /task/
		r.Put("/{taskID}", server.updateTask)    //PUT /task/123 - edit task
		r.Delete("/{taskID}", server.deleteTask) //DELETE /task/123 - delete dask
	})

	server.Router = r
	return server, nil
}
