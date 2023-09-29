package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"movie.api/internal/data"
	"movie.api/internal/validator"
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

// rateLimit restricts number of requests per second to server.
func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// background clean up goroutine to avoid rate limiter map overgrowth for clients' ip addresses.
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			// clean rate limiter map if client last seen was 3 minutes ago.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr) // clients ip extraction.
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(
						rate.Limit(app.config.limiter.rps),
						app.config.limiter.burst,
					),
				}
			}

			clients[ip].lastSeen = time.Now() // update last seen time.
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			mu.Unlock() // defer will wait until all chain will be completed.
		}
		next.ServeHTTP(w, r)
	})
	// the approach works only for single instance on a single-machine.
	// the different approach is required for distributed applications system.
	// different approach:
	// => 1. Redis usage.
	// => 2. Load balancer/reverse proxy =>  In-build rate limiter of Nginx.
}

// authenticate
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// header indicator for caches and indicates that response depends on Authorization header value.
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		// If empty, add anonymous user to context.
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		// Otherwise, validate and continue request chain processing.
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		token := headerParts[1]
		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
