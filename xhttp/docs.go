package xhttp

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"
	"path"
	"time"
)

//go:embed docs/web-components.min.js
var docsScript []byte

//go:embed docs/styles.min.css
var docsStyles []byte

//go:embed docs/index.gohtml
var docsTemplateSrc string

var docsTemplate = template.Must(template.New("").Parse(docsTemplateSrc))

// Returns http handler that provides OpenAPI specification page.
// Uses Stoplight Elements.
func DocumentationHandler(name, basePath string, spec []byte) http.Handler {
	var index bytes.Buffer
	if err := docsTemplate.Execute(&index, map[string]string{
		"Name":   name,
		"Prefix": path.Join("/", basePath),
	}); err != nil {
		panic("could not execute template: " + err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		switch r.URL.Path {
		case "/web-components.min.js":
			http.ServeContent(w, r, "/web-components.min.js", time.Time{}, bytes.NewReader(docsScript))
		case "/styles.min.css":
			http.ServeContent(w, r, "/styles.min.css", time.Time{}, bytes.NewReader(docsStyles))
		case "/openapi.yaml":
			http.ServeContent(w, r, "/openapi.yaml", time.Time{}, bytes.NewReader(spec))
		default:
			http.ServeContent(w, r, "/index.html", time.Time{}, bytes.NewReader(index.Bytes()))
		}
	})
}

// Adds the given prefix to the documentation handler
// and returns pattern and http handler for (*http.ServeMux).Handle-like methods.
func Documentation(name, basePath, prefix string, spec []byte) (string, http.Handler) {
	handler := DocumentationHandler(name, path.Join(basePath, prefix), spec)
	return Prefix(prefix, handler)
}
