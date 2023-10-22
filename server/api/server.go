package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	db "github.com/punkzberryz/todo/db/sqlc"
	"github.com/punkzberryz/todo/token"
	"github.com/punkzberryz/todo/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	Router     *chi.Mux
}

// Create new HTTP server and setup routing
func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		config:     config,
		store:      *store,
		tokenMaker: tokenMaker,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	//user-route
	r.Route("/user", func(r chi.Router) {
		r.Post("/", server.createUser)     //POST /user/
		r.Post("/login", server.loginUser) //POST /user/login
	})
	// user-route-protected
	r.Route("/me", func(r chi.Router) {
		r.Use(server.authMiddleware)
		r.Get("/", server.getCurrentUser) //GET /me/
	})
	//task-route
	r.Route("/task", func(r chi.Router) {
		r.Use(server.authMiddleware)       //require Header {Authorization: Bearer token}
		r.Get("/{taskID}", server.getTask) //GET /task/123
		r.Post("/", server.createTask)     //POST /task/123
		r.Get("/", server.getTaskList)     //GET /task/
	})

	server.Router = r
	return server, nil
}
