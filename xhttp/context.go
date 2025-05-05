package xhttp

import (
	"context"
	"net/http"
)

type contextKey struct{}

type contextData struct {
	err    error
	values map[string]any
}

type dataCtx struct {
	context.Context
	data *contextData
}

func (c *dataCtx) Value(key any) any {
	switch k := key.(type) {
	case contextKey:
		return c.data
	case string:
		if v, ok := c.data.values[k]; ok {
			return v
		}
	}
	return c.Context.Value(key)
}

// Middleware that provides a custom context that
// can be used to set and get values and errors
// inside handlers and other middlewares.
func Context() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(&dataCtx{
				Context: r.Context(),
				data: &contextData{
					err:    nil,
					values: make(map[string]any),
				},
			}))
		})
	}
}

// Puts key-value pair in the context provided by Context middleware.
// NOTE: use Get to retrieve the value, using (context.Context).Value won't work.
func Set[T any](r *http.Request, key string, value T) {
	data, ok := r.Context().Value(contextKey{}).(*contextData)
	if !ok {
		panic("invalid context")
	}
	data.values[key] = value
}

// Looks up a key from the context provided by Context middleware.
func Get[T any](r *http.Request, key string) (val T, ok bool) {
	data, ok := r.Context().Value(contextKey{}).(*contextData)
	if !ok {
		panic("invalid context")
	}

	v, ok := data.values[key]
	if !ok {
		return val, false
	}

	val, ok = v.(T)
	return val, ok
}

// Returns an error saved in the context provided by Context middleware.
func GetError(r *http.Request) error {
	data, ok := r.Context().Value(contextKey{}).(*contextData)
	if !ok {
		panic("invalid context")
	}
	return data.err
}

// Saves an error in the context provided by Context middleware.
func Error(r *http.Request, err error) {
	data, ok := r.Context().Value(contextKey{}).(*contextData)
	if !ok {
		panic("invalid context")
	}
	data.err = err
}
