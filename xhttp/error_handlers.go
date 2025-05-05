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

// Middleware that allows to catch http.StatusNotFound codes
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

// Middleware that allows to catch http.StatusMethodNotAllowed codes
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
