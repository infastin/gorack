package xrest

import (
	"net/http"
	"net/http/pprof"
	"strings"
)

func Debug(prefix string) http.Handler {
	return http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
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
	}))
}
