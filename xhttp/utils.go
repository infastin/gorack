package xhttp

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"
	"unsafe"
)

func toBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func routePathJoin(elems ...string) string {
	elems = append([]string{"/"}, elems...)
	result := path.Join(elems...)
	if result[len(result)-1] != '/' {
		last := elems[len(elems)-1]
		if last != "" && last[len(last)-1] == '/' {
			result += "/"
		}
	}
	return result
}

func muxPrefix(prefix string) (pattern, strip string) {
	pattern = routePathJoin("/", prefix)
	if pattern[len(pattern)-1] != '/' {
		pattern += "/"
	}
	return pattern, pattern[:len(pattern)-1]
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
