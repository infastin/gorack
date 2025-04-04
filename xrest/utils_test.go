package xrest

import (
	"testing"
)

func TestRoutePathJoin(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{input: []string{}, expected: "/"},
		{input: []string{""}, expected: "/"},
		{input: []string{"/"}, expected: "/"},
		{input: []string{"//"}, expected: "/"},
		{input: []string{"/v1"}, expected: "/v1"},
		{input: []string{"/v1/"}, expected: "/v1/"},
		{input: []string{"/v1//"}, expected: "/v1/"},
		{input: []string{"/v1", "foo", "bar"}, expected: "/v1/foo/bar"},
		{input: []string{"/v1", "/foo", "/bar"}, expected: "/v1/foo/bar"},
		{input: []string{"/v1/", "/foo/", "/bar"}, expected: "/v1/foo/bar"},
		{input: []string{"/v1/", "/foo/", "/bar/"}, expected: "/v1/foo/bar/"},
		{input: []string{"/v1/", "/foo/", "bar/"}, expected: "/v1/foo/bar/"},
	}
	for _, tt := range tests {
		output := routePathJoin(tt.input...)
		if output != tt.expected {
			t.Errorf("routePathJoin mismatch: input=%q expected=%q got=%q",
				tt.input, tt.expected, output)
			return
		}
	}
}

func TestMuxPrefix(t *testing.T) {
	tests := []struct {
		input           string
		expectedPattern string
		expectedStrip   string
	}{
		{input: "", expectedPattern: "/", expectedStrip: ""},
		{input: "/", expectedPattern: "/", expectedStrip: ""},
		{input: "/v1", expectedPattern: "/v1/", expectedStrip: "/v1"},
		{input: "/v1/", expectedPattern: "/v1/", expectedStrip: "/v1"},
		{input: "/v1/foo/bar", expectedPattern: "/v1/foo/bar/", expectedStrip: "/v1/foo/bar"},
		{input: "/v1/foo/bar/", expectedPattern: "/v1/foo/bar/", expectedStrip: "/v1/foo/bar"},
	}
	for _, tt := range tests {
		pattern, strip := muxPrefix(tt.input)
		if pattern != tt.expectedPattern {
			t.Errorf("muxPrefix pattern mismatch: input=%q expected=%q got=%q",
				tt.input, tt.expectedPattern, pattern)
			return
		}
		if strip != tt.expectedStrip {
			t.Errorf("muxPrefix strip mismatch: input=%q expected=%q got=%q",
				tt.input, tt.expectedStrip, strip)
			return
		}
	}
}

func TestFullExt(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "foo", expected: ""},
		{input: "foo.mp3", expected: ".mp3"},
		{input: "foo.", expected: "."},
		{input: "foo.png", expected: ".png"},
		{input: "index.html.tmpl", expected: ".html.tmpl"},
		{input: "hello.tar.gz", expected: ".tar.gz"},
	}
	for _, tt := range tests {
		output := fullExt(tt.input)
		if output != tt.expected {
			t.Errorf("fullExt mismatch: input=%q expected=%q got=%q",
				tt.input, tt.expected, output)
			return
		}
	}
}
