package xrest

import (
	"net/url"
)

func muxPrefix(prefix string) (pattern, strip string, err error) {
	pattern, err = url.JoinPath("/", prefix, "/")
	if err != nil {
		return "", "", err
	}
	return pattern, pattern[:len(pattern)-1], nil
}
