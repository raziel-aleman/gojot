package gojot

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiterConfig holds the configuration for the rate limiter.
type RateLimiterConfig struct {
	RequestLimit int
	TimeWindow   int // In seconds
}

// RequestLog store requests timestamps for each user/IP.
type RequestLog struct {
	mu       sync.Mutex // Mutex for concurrency safety
	requests map[string][]time.Time
	config   RateLimiterConfig
}

var requestLog RequestLog

// SetRateLimiterConfig initializes the rate limiter with the given requests limit and time window (in seconds).
// Not initializing the rate limiter will casue a panic due to a nil requestLog store struct.
// Default configuration to be implemented.
func SetRateLimiterConfig(requestLimit, timeWindow int) {
	requestLog = RequestLog{
		requests: make(map[string][]time.Time),
		config:   RateLimiterConfig{RequestLimit: requestLimit, TimeWindow: timeWindow},
	}
}

// isRequestAllowed checks if the request is allowed with the given requests limit and time window provided when
// initializing rate limiter.
func isRequestAllowed(userID string) bool {
	requestLog.mu.Lock()
	defer requestLog.mu.Unlock()

	currentTime := time.Now()

	userRequests := requestLog.requests[userID]

	// Remove outdated requests
	var validRequests []time.Time
	for _, timestamp := range userRequests {
		if timestamp.After(currentTime.Add(-time.Duration(requestLog.config.TimeWindow) * time.Second)) {
			validRequests = append(validRequests, timestamp)
		}
	}

	requestLog.requests[userID] = validRequests

	if len(validRequests) < requestLog.config.RequestLimit {
		requestLog.requests[userID] = append(validRequests, currentTime)
		return true
	}

	return false
}

// RateLimiterMiddleware checks if the request is allowed.
// If the request is not allowed, it returns a too many requests status code.
func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		clientIP := r.RemoteAddr // Use remote IP address for simplicity

		if isRequestAllowed(clientIP) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	})
}
