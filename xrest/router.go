package xrest

import (
	"net/http"
	"strings"
)

type Router struct {
	mux        *http.ServeMux
	prefix     string
	strip      string
	middleware Middleware
}

func NewRouter() *Router {
	return &Router{
		mux:        http.NewServeMux(),
		prefix:     "",
		strip:      "",
		middleware: nil,
	}
}

func (r *Router) Group(prefix string, fn func(*Router)) {
	prefix, strip := muxPrefix(prefix)
	g := &Router{
		mux:        r.mux,
		prefix:     prefix,
		strip:      strip,
		middleware: r.middleware,
	}
	fn(g)
}

func (r *Router) Use(middlewares ...Middleware) {
	if len(middlewares) == 0 {
		return
	}
	if r.middleware != nil {
		middlewares = append([]Middleware{r.middleware}, middlewares...)
	}
	r.middleware = Chain(middlewares...)
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	if r.prefix != "" {
		slashIdx := strings.IndexByte(pattern, '/')
		if slashIdx == -1 {
			panic("xrest: invalid pattern")
		}
		// [METHOD ][HOST]<PREFIX>[/PATH]
		pattern = pattern[:slashIdx] + routePathJoin(r.prefix, pattern[slashIdx:])
	}
	if r.strip != "" {
		handler = http.StripPrefix(r.strip, handler)
	}
	if r.middleware != nil {
		r.mux.Handle(pattern, r.middleware(handler))
	} else {
		r.mux.Handle(pattern, handler)
	}
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.Handle(pattern, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
