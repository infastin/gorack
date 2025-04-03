package xrest

import "net/http"

type notFoundWriter struct {
	http.ResponseWriter
	request    *http.Request
	statusCode int
	handler    http.HandlerFunc
}

func (w *notFoundWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusNotFound {
		w.statusCode = statusCode
		return
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *notFoundWriter) Write(data []byte) (int, error) {
	if w.statusCode == http.StatusNotFound {
		w.handler(w.ResponseWriter, w.request)
		return 0, nil
	}
	return w.ResponseWriter.Write(data)
}

// Middleware that allows to catch http.StatusNotFound codes
// written to http.ResponseWriter and handle them.
func NotFound(handler http.HandlerFunc) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nw := &notFoundWriter{
				ResponseWriter: w,
				request:        r,
				statusCode:     http.StatusOK,
				handler:        handler,
			}
			next.ServeHTTP(nw, r)
		})
	}
}
