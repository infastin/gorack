package xrest

import (
	"path"
	"strings"
)

func muxPrefix(prefix string) (pattern, strip string) {
	pattern = path.Clean(prefix) + "/"
	strip = strings.TrimRight(pattern, "/")
	return pattern, strip
}
