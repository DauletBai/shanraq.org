package main

import (
	"context"
	"net/http"

	"myblog/blog"
)

// userKey is a type of our own, and that is the whole point: a value in a
// context is found by a key compared together with its type. A string "user"
// from our package and a string "user" from somebody's library are the same
// key, and whoever writes last wins. A type declared here cannot be repeated
// from outside.
type userKey struct{}

// withUser looks the signed-in person up once per request and puts what it
// found into the request's context. Before this every handler asked the
// database for itself.
func (app *app) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(blog.CookieName); err == nil {
			if u, err := app.store.UserByToken(r.Context(), c.Value); err == nil {
				// WithContext returns a copy: a request is not changed in
				// place, so the result has to be assigned back — the same trap
				// as append in the slices lesson.
				r = r.WithContext(context.WithValue(r.Context(), userKey{}, u))
			}
		}
		next.ServeHTTP(w, r)
	})
}

// userFrom hides the key: outside this file nobody needs to know that the
// answer travels in a context at all.
func userFrom(ctx context.Context) (blog.User, bool) {
	u, ok := ctx.Value(userKey{}).(blog.User)
	return u, ok
}
