package main

import (
	"fmt"
	"net/http"
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
