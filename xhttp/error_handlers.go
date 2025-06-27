package xhttp

import "net/http"

type statusCatcher struct {
	http.ResponseWriter
	request   *http.Request
	catchCode int
	handler   http.HandlerFunc
	didCatch  bool
	didHandle bool
}

func (w *statusCatcher) WriteHeader(statusCode int) {
	if statusCode == w.catchCode {
		w.didCatch = true
		return
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusCatcher) Write(data []byte) (int, error) {
	if w.didCatch {
		if !w.didHandle {
			w.handler(w.ResponseWriter, w.request)
			w.didHandle = true
		}
		return 0, nil
	}
	return w.ResponseWriter.Write(data)
}

// NotFound is a middleware that allows to catch http.StatusNotFound codes
// written to http.ResponseWriter and handle them.
func NotFound(handler http.HandlerFunc) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scw := &statusCatcher{
				ResponseWriter: w,
				request:        r,
				catchCode:      http.StatusNotFound,
				handler:        handler,
				didCatch:       false,
				didHandle:      false,
			}
			next.ServeHTTP(scw, r)
		})
	}
}

// MethodNotAllowed is a middleware that allows to catch http.StatusMethodNotAllowed codes
// written to http.ResponseWriter and handle them.
func MethodNotAllowed(handler http.HandlerFunc) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scw := &statusCatcher{
				ResponseWriter: w,
				request:        r,
				catchCode:      http.StatusMethodNotAllowed,
				handler:        handler,
				didCatch:       false,
				didHandle:      false,
			}
			next.ServeHTTP(scw, r)
		})
	}
}

// RemoveErrorHandlers is a middleware that removes wrappers around the original ResponseWriter
// that were added by NotFound and MethodNotAllowed middlewares.
func RemoveErrorHandlers(handler http.HandlerFunc) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for {
				scw, ok := w.(*statusCatcher)
				if !ok {
					break
				}
				w = scw.ResponseWriter
			}
			next.ServeHTTP(w, r)
		})
	}
}
