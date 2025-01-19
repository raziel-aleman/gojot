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

	"github.com/go-chi/chi/v5"
	"github.com/raziel-aleman/gojot"
)

func main() {
	r := chi.NewRouter()

	// Add the optional middlewares for the chi router
	gojot.HelperMiddlewares(r)

	// Initialize rate limiter config
	gojot.SetRateLimiterConfig(100, 60)

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Use(gojot.RateLimiterMiddleware)
		r.Get("/login", LoginHandler)
	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Use(gojot.RateLimiterMiddleware)
		r.Use(gojot.AuthMiddleware([]byte("your-secret-key")))
		r.Get("/protected", ProtectedHandler)
		r.Get("/logout", LogoutHandler)
	})

	go func() {
		if err := http.ListenAndServe(":8080", r); err != http.ErrServerClosed {
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
