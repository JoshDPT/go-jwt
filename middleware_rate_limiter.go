package main

import (
	"net/http"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(2, 5) // Allowing 2 requests per 5 seconds

// rateLimitMiddleware is a middleware function that performs rate limiting.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
