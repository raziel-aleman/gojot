package main

import (
	"context"
	"fmt"
	"io"
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

	// Use the JWT helper middlewares
	gojot.HelperMiddlewares(r)

	// Protected route
	r.Get("/protected", gojot.JotMiddleware([]byte("your-secret-key"), func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(gojot.ContextKeyUserId).(uint64)
		fmt.Fprintf(w, "Welcome, user %d!", userID)
	}))

	// Route to generate and add token to a request
	r.Get("/generate-token", func(w http.ResponseWriter, r *http.Request) {
		// Assuming you have a user ID from authentication
		userID := uint64(123456)

		// Generate the JWT token
		token, err := gojot.GenerateToken(userID, []byte("your-secret-key"), time.Hour*24)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		/*
			// Creat JWT cookie for browser

			cookie := &http.Cookie{
				Name:     "jwt-token", // <- should be any unique key you want
				Value:    token,       // <- the token, recommend to encode by SecureCookie
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
		*/

		// Create a new request with the token in the Authorization header
		req, err := http.NewRequest("GET", "http://localhost:8080/protected", nil)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Authorization", "Bearer "+token)

		// Send the request to the protected endpoint
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to send request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Handle the response
		body, _ := io.ReadAll(resp.Body)
		//fmt.Println(w, string(body))
		w.Write(body)
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
