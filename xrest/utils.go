package xrest

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func muxPrefix(prefix string) (pattern, strip string, err error) {
	pattern, err = url.JoinPath("/", prefix, "/")
	if err != nil {
		return "", "", err
	}
	return pattern, pattern[:len(pattern)-1], nil
}

func fullExt(filename string) string {
	filename = path.Base(filename)
	if pos := strings.IndexByte(filename, '.'); pos != -1 {
		return filename[pos:]
	}
	return ""
}

func detectContentType(filename string, content io.Reader) (mimetype string, reader io.Reader) {
	mimetype = mime.TypeByExtension(fullExt(filename))
	if mimetype == "" || mimetype == "application/octet-stream" {
		var sniff [512]byte
		n, _ := io.ReadFull(content, sniff[:])
		mimetype = http.DetectContentType(sniff[:n])
		content = io.MultiReader(bytes.NewReader(sniff[:n]), content)
	}
	return mimetype, content
}
