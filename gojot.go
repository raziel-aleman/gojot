package gojot

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type contextKey string

var ContextKeyUserId = contextKey("userId")
var ContextKeyToken = contextKey("token")

// JotMiddleware is a middleware function that validates JWT tokens
func JotMiddleware(secretKey []byte, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextKeyToken, r.Header.Get("Authorization"))
		claims, err := ValidateToken(ctx, secretKey)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, ContextKeyUserId, claims.UserID)
		r = r.WithContext(ctx)
		handlerFunc(w, r)
	}
}

// UseMiddleware is a helper function to use the middleware with Chi router
func HelperMiddlewares(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}
