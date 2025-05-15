package xhttp_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/infastin/gorack/xhttp"
)

func ExampleBindParams() {
	mux := http.NewServeMux()

	type Anonymous struct {
		A string `query:"a"`
		B string `query:"b"`
	}

	type Inline struct {
		C string `query:"c"`
		D string `query:"d"`
	}

	type Params struct {
		Anonymous
		Inline      Inline `inline:""`
		Foo         int    `query:"foo"`
		NullableFoo *int   `query:"nullable_foo"`
		Bar         int    `header:"bar"`
		Baz         string `path:"baz"`
	}

	mux.HandleFunc("GET /params/{baz}", func(w http.ResponseWriter, r *http.Request) {
		var params Params
		if err := xhttp.BindParams(r, &params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		xhttp.JSON(w, http.StatusOK, &params)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	uri, _ := url.Parse(ts.URL)
	query := url.Values{
		"a":   []string{"I am A"},
		"b":   []string{"Hello from B"},
		"c":   []string{"Hello from C"},
		"d":   []string{"I am D"},
		"foo": []string{"42"},
	}
	uri.Path = "/params/qux"
	uri.RawQuery = query.Encode()

	resp, err := ts.Client().Do(&http.Request{
		Method: http.MethodGet,
		URL:    uri,
		Header: http.Header{"bar": []string{"123"}},
	})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", data)

	// Output: {"A":"I am A","B":"Hello from B","Inline":{"C":"Hello from C","D":"I am D"},"Foo":42,"NullableFoo":null,"Bar":123,"Baz":"qux"}
}
