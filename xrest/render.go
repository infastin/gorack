package xrest

import (
	"encoding/json"
	htmltemplate "html/template"
	"io"
	"mime"
	"net/http"
	"path"
	"strconv"
	texttemplate "text/template"
)

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Text(w http.ResponseWriter, code int, body []byte) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(code)
	w.Write(body) //nolint:errcheck
}

func TextTemplate(w http.ResponseWriter, code int, tmpl *texttemplate.Template, data any) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	tmpl.Execute(w, data) //nolint:errcheck
}

func HTML(w http.ResponseWriter, code int, body []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(code)
	w.Write(body) //nolint:errcheck
}

func HTMLTemplate(w http.ResponseWriter, code int, tmpl *htmltemplate.Template, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	tmpl.Execute(w, data) //nolint:errcheck
}

func JSON(w http.ResponseWriter, code int, body any) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	enc.Encode(body) //nolint:errcheck
}

func File(w http.ResponseWriter, code int, filename string, content io.Reader, size int64) {
	mimetype := mime.TypeByExtension(path.Ext(filename))
	if mimetype == "" {
		mimetype = "application/octet-stream"
	}

	w.Header().Set("Content-Type", mimetype)
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.WriteHeader(code)

	io.Copy(w, content) //nolint:errcheck
}

func StreamFile(w http.ResponseWriter, code int, filename string, content io.Reader) {
	mimetype := mime.TypeByExtension(path.Ext(filename))
	if mimetype == "" {
		mimetype = "application/octet-stream"
	}

	w.Header().Set("Content-Type", mimetype)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.WriteHeader(code)

	io.Copy(w, content) //nolint:errcheck
}

func Data(w http.ResponseWriter, code int, data []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(code)
	w.Write(data) //nolint:errcheck
}

func Stream(w http.ResponseWriter, code int, content io.Reader) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(code)
	io.Copy(w, content) //nolint:errcheck
}
