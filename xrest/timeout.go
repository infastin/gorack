package xrest

import (
	goctx "context"
	"net/http"
	"time"
)

func Timeout(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := goctx.WithTimeout(r.Context(), timeout)
			defer cancel()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func TimeoutCause(timeout time.Duration, cause error) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := goctx.WithTimeoutCause(r.Context(), timeout, cause)
			defer cancel()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
