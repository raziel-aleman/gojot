package gojot

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddlewareNoCookie(t *testing.T) {
	// Create a mock handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a middleware instance
	authMiddleware := AuthMiddleware([]byte("your-secret-key"))

	// Create a request and recorder
	req := httptest.NewRequest("GET", "/protected", nil)
	res := httptest.NewRecorder()

	// Execute the middleware
	authMiddleware(nextHandler).ServeHTTP(res, req)

	// Check the response status code
	if res.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
	}
}

func TestAuthMiddlewareWithCookie(t *testing.T) {
	// Create a mock handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a middleware instance
	authMiddleware := AuthMiddleware([]byte("your-secret-key"))

	// Create a request and recorder
	req := httptest.NewRequest("GET", "/protected", nil)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJfaWQiOjEyMzQ1NiwibmFtZSI6IlRlc3QiLCJlbWFpbCI6InRlc3RAZW1haWwuY29tIn0sImV4cCI6MTczNDUwNjQwN30.hqyD6pim7lLedGw-DaoVKuWD458T8Nc83DNDQbiezuM"
	req.AddCookie(&http.Cookie{
		Name:     authCookieName, // <- should be any unique key you want
		Value:    token,          // <- the token, recommend to encode by SecureCookie
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(time.Hour * 24)),
	})
	res := httptest.NewRecorder()

	// Execute the middleware
	authMiddleware(nextHandler).ServeHTTP(res, req)

	// Check the response status code
	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
	}
}
