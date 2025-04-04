package xrest

import "net/http"

// Adds the given prefix to the given http handler
// and returns pattern and http handler for (*http.ServeMux).Handle-like methods.
func Prefix(prefix string, handler http.Handler) (string, http.Handler) {
	pattern, strip := muxPrefix(prefix)
	return pattern, http.StripPrefix(strip, handler)
}
