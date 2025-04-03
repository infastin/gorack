package xrest

import (
	"net/http"
)

// Adds the given prefix to the http handler returned from fn
// and returns pattern and http handler for (*http.ServeMux).Handle-like methods.
func Group(prefix string, fn func() http.Handler) (string, http.Handler) {
	return Prefix(prefix, fn())
}
