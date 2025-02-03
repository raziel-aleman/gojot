package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raziel-aleman/gojot"
)

func main() {

	// Standard library
	mux := http.NewServeMux()
	mux.Handle("GET /login",
		gojot.RateLimiterMiddleware(http.HandlerFunc(LoginHandler)))
	mux.Handle("GET /protected",
		gojot.RateLimiterMiddleware(gojot.AuthMiddleware([]byte("your-secret-key"))(http.HandlerFunc(ProtectedHandler))))
	mux.Handle("GET /logout",
		gojot.RateLimiterMiddleware(gojot.AuthMiddleware([]byte("your-secret-key"))(http.HandlerFunc(LogoutHandler))))

	// Initialize rate limiter config
	gojot.SetRateLimiterConfig(100, 60)

	go func() {
		if err := http.ListenAndServe(":8080", mux); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %s", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-interrupt
	log.Println("Shutting down server...")

	gracefulCtx, gracefulCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer gracefulCancel()

	if err := (*http.Server).Shutdown(&http.Server{}, gracefulCtx); err != nil {
		log.Printf("Shutdown error: %v\n", err)
		defer os.Exit(1)
		return
	} else {
		log.Printf("Server stopped\n")
	}

	defer os.Exit(0)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(gojot.ContextKeyUser).(gojot.User)
	log.Println(user)
	fmt.Fprintf(w, "Welcome, user %d!", user.ID)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Assuming you have a user ID from authentication
	user := gojot.User{ID: uint64(123456), Name: "Test", Email: "test@email.com"}

	// Generate the JWT token
	token, err := gojot.GenerateToken(user, []byte("your-secret-key"), time.Hour*24)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set auth cookie with JWT token
	gojot.SetAuthCookie(w, token, time.Hour*24)

	// Redirect to protected endpoint
	http.Redirect(w, r, "http://localhost:8080/protected", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Remove auth cookie with JWT token
	gojot.RemoveAuthCookie(w)

	// Redirect to protected endpoint
	http.Redirect(w, r, "http://localhost:8080/protected", http.StatusFound)
}
