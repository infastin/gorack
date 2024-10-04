package xrest

import (
	"net/http"
)

func Group(prefix string, fn func() http.Handler) (string, http.Handler) {
	pattern, strip, err := muxPrefix(prefix)
	if err != nil {
		panic("invalid prefix: " + err.Error())
	}
	return pattern, http.StripPrefix(strip, fn())
}
