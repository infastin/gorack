package xhttp

import (
	"context"
	"net/http"
	"time"
)

// Timeout is a middleware that will cancel request's context after the specified duration.
func Timeout(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// TimeoutCause is a middleware that will cancel request's context with the given cause
// after the specified duration
func TimeoutCause(timeout time.Duration, cause error) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeoutCause(r.Context(), timeout, cause)
			defer cancel()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
