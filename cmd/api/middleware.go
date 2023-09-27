package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

// recoverPanic recovers work and send error response.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Trigger to make Go's HTTP server automatically close the current connection after a response has been sent.
				w.Header().Set("Connection", "close")
				// recover value = any => string => log it.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// Initialize a new rate limiter which allows an average of 2 requests per second, with a maximum of 4 requests in a single ‘burst’.
	limiter := rate.NewLimiter(2, 4)
	// The function we are returning is a closure, which 'closes over' the limiter variable.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call limiter.Allow() to see if the request is permitted, and if it's not,
		// then rateLimitExceededResponse() returns a 429 Too Many Requests response.
		if !limiter.Allow() { // Allow takes one token from bucket.
			app.rateLimitExceededResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
