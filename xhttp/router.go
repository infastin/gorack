package xhttp

import (
	"net/http"
	"strings"
)

// Router is simple wrapper around http.ServerMux that introduces
// additional convenient methods.
type Router struct {
	mux         *http.ServeMux
	prefix      string
	strip       string
	middlewares []Middleware
}

// NewRouter creates a new Router.
func NewRouter() *Router {
	return &Router{
		mux:         http.NewServeMux(),
		prefix:      "",
		strip:       "",
		middlewares: make([]Middleware, 0),
	}
}

// Group creates a new router group with prefix.
func (r *Router) Group(prefix string, fn func(group *Router)) {
	prefix, strip := muxPrefix(prefix)
	g := &Router{
		mux:         r.mux,
		prefix:      prefix,
		strip:       strip,
		middlewares: r.middlewares,
	}
	fn(g)
}

// Mount attaches handler as a subrouter with prefix.
func (r *Router) Mount(prefix string, handler http.Handler) {
	r.Handle(Prefix(prefix, handler))
}

// Use adds middlewares to the chain.
func (r *Router) Use(middlewares ...Middleware) {
	if len(middlewares) == 0 {
		return
	}
	r.middlewares = append(r.middlewares, middlewares...)
}

// Handle registers the handler for the given pattern.
func (r *Router) Handle(pattern string, handler http.Handler) {
	if r.prefix != "" {
		slashIdx := strings.IndexByte(pattern, '/')
		if slashIdx == -1 {
			panic("xhttp: invalid pattern")
		}
		// [METHOD ][HOST]<PREFIX>[/PATH]
		pattern = pattern[:slashIdx] + routePathJoin(r.prefix, pattern[slashIdx:])
	}
	if r.strip != "" {
		handler = http.StripPrefix(r.strip, handler)
	}
	if len(r.middlewares) != 0 {
		wrapper := Chain(r.middlewares...)
		r.mux.Handle(pattern, wrapper(handler))
	} else {
		r.mux.Handle(pattern, handler)
	}
}

// HandleFunc registers the handler function for the given pattern.
func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.Handle(pattern, handler)
}

// ServeHTTP implements http.Handler interface, which serves HTTP requests.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
