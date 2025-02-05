package gojot

import (
	"context"
	"net/http"
	"time"
)

type contextKey string

const ContextKeyUser = contextKey("user")
const authCookieName = "access_token"

// AuthMiddleware validates checks for JWT tokens and validates them.
// If request has a valid JWT token, the user information is added to the request context,
// otherwise an unauthorized status code is returned.
func AuthMiddleware(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := ExtractTokenFromCookie(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := ValidateToken(token, secretKey)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUser, claims.User)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// // HelperMiddlewares adds common middlewares (RequestID, RealIP, Logger, Recoverer) to a Chi router.
// func HelperMiddlewares(r chi.Router) {
// 	r.Use(middleware.RequestID)
// 	r.Use(middleware.RealIP)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)
// }

// SetAuthCookie sets auth cookie with JWT token for client.
func SetAuthCookie(w http.ResponseWriter, token string, expirationTime time.Duration) {
	cookie := &http.Cookie{
		Name:     authCookieName, // <- should be any unique key you want
		Value:    token,          // <- the token, recommend to encode by SecureCookie
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(expirationTime),
	}
	http.SetCookie(w, cookie)
}

// RemoveAuthCookie deletes auth cookie with JWT token from client.
func RemoveAuthCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}

// ExtractTokenFromCookie extracts JWT token from requests auth cookie.
func ExtractTokenFromCookie(r *http.Request) (string, error) {
	jwtCookie, err := r.Cookie(authCookieName)
	if err != nil {
		return "", err
	}

	return jwtCookie.Value, nil
}
