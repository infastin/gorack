package xrest

import (
	"net/http"
	"net/http/pprof"
	"strings"
)

// Returns http handler that provides /pprof routes.
func ProfilerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		name, found := strings.CutPrefix(r.URL.Path, "/pprof/")
		if !found {
			http.NotFound(w, r)
			return
		}

		switch r.URL.Path {
		case "/pprof/cmdline":
			pprof.Cmdline(w, r)
		case "/pprof/profile":
			pprof.Profile(w, r)
		case "/pprof/symbol":
			pprof.Symbol(w, r)
		case "/pprof/trace":
			pprof.Trace(w, r)
		default:
			if name != "" {
				pprof.Handler(name).ServeHTTP(w, r)
			} else {
				pprof.Index(w, r)
			}
		}
	})
}

// Adds the given prefix to the profiler handler
// and returns pattern and http handler for (*http.ServeMux).Handle-like methods.
func Profiler(prefix string) (string, http.Handler) {
	return Prefix(prefix, ProfilerHandler())
}

// Deprecated: use Profiler.
func Debug(prefix string) (string, http.Handler) {
	return Profiler(prefix)
}
