package main

import (
	"context"
	"net/http"

	"movie.api/internal/data"
)

// contextKey
type contextKey string

// constant key needed to get/set user details in request context.
const userContextKey = contextKey("user")

// contextSetUser returns a new copy of the request context with User struct added.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrieves the User struct from the request context. The only
// time that it would be used when it is logically expected there to be User struct
// value in the context, and if it doesn't exist it will firmly be an 'unexpected' error.
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
