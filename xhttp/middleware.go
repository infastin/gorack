package xhttp

import "net/http"

type Middleware func(next http.Handler) http.Handler

// Returns a new middleware that is the result of
// chaining multiple middlewares.
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		h := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}
