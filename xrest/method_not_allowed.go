package xrest

import "net/http"

type methodNotAllowedWriter struct {
	http.ResponseWriter
	request    *http.Request
	statusCode int
	handler    http.HandlerFunc
}

func (w *methodNotAllowedWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusMethodNotAllowed {
		w.statusCode = statusCode
		return
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *methodNotAllowedWriter) Write(data []byte) (int, error) {
	if w.statusCode == http.StatusMethodNotAllowed {
		w.handler(w.ResponseWriter, w.request)
		return 0, nil
	}
	return w.ResponseWriter.Write(data)
}

func MethodNotAllowed(handler http.HandlerFunc) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw := &methodNotAllowedWriter{
				ResponseWriter: w,
				request:        r,
				statusCode:     http.StatusOK,
				handler:        handler,
			}
			next.ServeHTTP(mw, r)
		})
	}
}
