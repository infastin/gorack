package xrest

import (
	goctx "context"
	"net/http"
)

type context struct {
	goctx.Context
	err    error
	values map[string]any
}

func (c *context) Value(key any) any {
	if k, ok := key.(string); ok {
		if v, ok := c.values[k]; ok {
			return v
		}
	}
	return c.Context.Value(key)
}

func Context() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(&context{
				Context: r.Context(),
				err:     nil,
				values:  make(map[string]any),
			}))
		})
	}
}

func Set[T any](r *http.Request, key string, value T) {
	c, ok := r.Context().(*context)
	if !ok {
		panic("invalid context")
	}
	c.values[key] = value
}

func Get[T any](r *http.Request, key string) (val T, ok bool) {
	c, ok := r.Context().(*context)
	if !ok {
		panic("invalid context")
	}

	v, ok := c.values[key]
	if !ok {
		return val, false
	}

	val, ok = v.(T)
	return val, ok
}

func GetError(r *http.Request) error {
	c, ok := r.Context().(*context)
	if !ok {
		panic("invalid context")
	}
	return c.err
}

func Error(r *http.Request, err error) {
	c, ok := r.Context().(*context)
	if !ok {
		panic("invalid context")
	}
	c.err = err
}
