package xhttp

import (
	"encoding/json"
	"encoding/xml"
	htmltemplate "html/template"
	"io"
	"net/http"
	"strconv"
	texttemplate "text/template"

	"github.com/infastin/gorack/fastconv"
)

// Writes http headers with the given status code.
func NoContent(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

// Writes body with the given status code
// and Content-Type set to "text/plain; charset=utf-8".
// Also sets Content-Length to the size of body.
func Text(w http.ResponseWriter, code int, body []byte) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(code)
	w.Write(body) //nolint:errcheck
}

// Executes text template tmpl with data
// and writes the output with the given status code
// and Content-Type set to "text/plain; charset=utf-8".
func TextTemplate(w http.ResponseWriter, code int, tmpl *texttemplate.Template, data any) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	tmpl.Execute(w, data) //nolint:errcheck
}

// Writes body with the given status code
// and Content-Type set to "text/html; charset=utf-8".
// Also sets Content-Length to the size of body.
func HTML(w http.ResponseWriter, code int, body []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(code)
	w.Write(body) //nolint:errcheck
}

// Executes html template tmpl with data
// and writes the output with the given status code
// and Content-Type set to "text/html; charset=utf-8".
func HTMLTemplate(w http.ResponseWriter, code int, tmpl *htmltemplate.Template, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	tmpl.Execute(w, data) //nolint:errcheck
}

// Encodes body as json and writes the output with the given status code
// and Content-Type set to "application/json".
func JSON(w http.ResponseWriter, code int, body any) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	enc.Encode(body) //nolint:errcheck
}

// Encodes body as xml with <?xml> header
// and writes the output with the given status code
// and Content-Type set to "application/xml; charset=utf-8".
func XMLWithHeader(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(code)

	w.Write(fastconv.Bytes(xml.Header)) //nolint:errcheck
	xml.NewEncoder(w).Encode(body)      //nolint:errcheck
}

// Encodes body as xml and writes the output with the given status code
// and Content-Type set to "application/xml; charset=utf-8".
func XML(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(code)
	xml.NewEncoder(w).Encode(body) //nolint:errcheck
}

// Writes content as an attachment with the given filename and with the given status code.
// If size is not zero, sets Content-Length.
// If Content-Type is not set, tries to determine it from the extension of filename
// and content itself, falling back to "application/octet-stream" if it is unable to determine a valid MIME type,
// and sets Content-Type to the resulting MIME type.
// NOTE: It is recommended to use http.ServeContent instead of this function.
func Attachment(w http.ResponseWriter, code int, filename string, content io.Reader, size int64) {
	if mimetype := w.Header().Get("Content-Type"); mimetype == "" {
		mimetype, content = detectContentType(filename, content)
		w.Header().Set("Content-Type", mimetype)
	}
	if size != 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.WriteHeader(code)
	io.Copy(w, content) //nolint:errcheck
}

// Writes content as is with the given status code.
// If size is not zero, sets Content-Length.
// If Content-Type is not set, tries to determine it from the extension of filename
// and content itself, falling back to "application/octet-stream" if it is unable to determine a valid MIME type,
// and sets Content-Type to the resulting MIME type.
// NOTE: It is recommended to use http.ServeContent instead of this function.
func File(w http.ResponseWriter, code int, filename string, content io.Reader, size int64) {
	if mimetype := w.Header().Get("Content-Type"); mimetype == "" {
		mimetype, content = detectContentType(filename, content)
		w.Header().Set("Content-Type", mimetype)
	}
	if size != 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	}
	w.WriteHeader(code)
	io.Copy(w, content) //nolint:errcheck
}

// Writes data with the given status code.
// and Content-Type set to "application/octet-stream".
// Also sets Content-Length to the length of data.
func Blob(w http.ResponseWriter, code int, data []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(code)
	w.Write(data) //nolint:errcheck
}

// Writes content as is with the given status code.
// and Content-Type set to "application/octet-stream".
func Stream(w http.ResponseWriter, code int, content io.Reader) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(code)
	io.Copy(w, content) //nolint:errcheck
}
