package gojot

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"
)

func TestRateLimiterMiddleware(t *testing.T) {

	// Create a mock handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Initialize rate limiter config
	SetRateLimiterConfig(100, 60)

	// Create a middleware instance
	rateLimiterMiddleware := RateLimiterMiddleware

	// Slice for goroutines response codes
	resCodes := []int{}

	var wg sync.WaitGroup

	// 101 requests as rate limiter default config only allows for 100 requests every 60 seconds
	for i := 0; i < 101; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Create a request and recorder
			req := httptest.NewRequest("GET", "/protected", nil)
			res := httptest.NewRecorder()

			// Execute the middleware
			rateLimiterMiddleware(nextHandler).ServeHTTP(res, req)

			resCodes = append(resCodes, res.Code)
		}()
	}

	wg.Wait()

	if !slices.Contains(resCodes, http.StatusTooManyRequests) {
		t.Errorf("Expected one status code %d, but all requests got %d", http.StatusTooManyRequests, http.StatusOK)
	}
}
