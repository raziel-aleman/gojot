package gojot

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type contextKey string

const ContextKeyUser = contextKey("user")
const authCookieName = "access_token"

// AuthMiddleware is a middleware handler that validates JWT tokens
func AuthMiddleware(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := extractTokenFromCookie(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := ValidateToken(token, secretKey)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(context.Background(), ContextKeyUser, claims.User)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// UseMiddleware is a helper function to use the middleware with Chi router
func HelperMiddlewares(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}

// Set auth cookie for client
func SetAuthCookie(w *http.ResponseWriter, token string, expirationTime time.Duration) {
	cookie := &http.Cookie{
		Name:     authCookieName, // <- should be any unique key you want
		Value:    token,          // <- the token, recommend to encode by SecureCookie
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(expirationTime),
	}
	http.SetCookie(*w, cookie)
}

// Remove auth cookie from client
func RemoveAuthCookie(w *http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
	http.SetCookie(*w, cookie)
}

// Extract JWT token from cookie
func extractTokenFromCookie(r *http.Request) (string, error) {
	jwtCookie, err := r.Cookie(authCookieName)
	if err != nil {
		return "", err
	}

	return jwtCookie.Value, nil
}
