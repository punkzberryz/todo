package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type ctxKey string

const (
	payloadKey ctxKey = "payload"
)

func (server *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("authorization")
			fields := strings.Fields(authHeader)
			if len(fields) < 2 {
				err := errors.New("invalid authorization header format")
				render.Render(w, r, ErrUnauthorized(err))
				return
			}
			authorizationType := strings.ToLower(fields[0])
			if authorizationType != "bearer" {
				err := fmt.Errorf("unsupported authorization type %s", authorizationType)
				render.Render(w, r, ErrUnauthorized(err))
				return
			}

			accessToken := fields[1]
			payload, err := server.token.Maker.VerifyToken(accessToken)
			if err != nil {
				render.Render(w, r, ErrUnauthorized(err))
				return
			}

			ctx := context.WithValue(r.Context(), payloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
}
