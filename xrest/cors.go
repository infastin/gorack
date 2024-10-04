package xrest

import (
	"net/http"
	"strings"
)

type CORSConfig struct {
	AllowOrigin  string
	AllowMethods []string
}

var DefaultConfig = CORSConfig{
	AllowOrigin: "*",
	AllowMethods: []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodPost,
		http.MethodDelete,
	},
}

func CORS() Middleware {
	return CORSWithConfig(DefaultConfig)
}

func CORSWithConfig(config CORSConfig) Middleware {
	if config.AllowOrigin == "" {
		config.AllowOrigin = DefaultConfig.AllowOrigin
	}
	if config.AllowMethods == nil {
		config.AllowMethods = DefaultConfig.AllowMethods
	}

	allowOrigin := config.AllowOrigin
	allowMethods := strings.Join(config.AllowMethods, ",")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Add("Access-Control-Allow-Methods", allowMethods)
			next.ServeHTTP(w, r)
		})
	}
}
