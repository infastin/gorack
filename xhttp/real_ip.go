package xhttp

import (
	"net"
	"net/http"
	"strings"
)

// RealIP extracts real ip from X-Forward-For and X-Real-IP headers.
func RealIP(r *http.Request) net.IP {
	ip := getIPFromXRealIP(r)
	if ip != nil {
		return ip
	}

	ip = getIPFromXForwardedFor(r)
	if ip != nil {
		return ip
	}

	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return net.ParseIP(host)
}

func getIPFromXForwardedFor(r *http.Request) net.IP {
	headers := r.Header.Values("X-Forwarded-For")
	if len(headers) == 0 {
		return nil
	}

	for _, header := range headers {
		parts := strings.Split(header, ",")
		for i := range parts {
			part := strings.TrimSpace(parts[i])
			ip := net.ParseIP(part)
			if ip != nil {
				return ip
			}
		}
	}

	return nil
}

func getIPFromXRealIP(r *http.Request) net.IP {
	xRealIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if xRealIP == "" {
		return nil
	}

	ip := net.ParseIP(xRealIP)
	if ip != nil {
		return ip
	}

	return nil
}
