package xrest

import (
	"net/http"
)

func Group(prefix string, fn func() http.Handler) (string, http.Handler) {
	pattern, prefix := muxPrefix(prefix)
	return pattern, http.StripPrefix(prefix, fn())
}
