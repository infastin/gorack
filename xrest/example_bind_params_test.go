package xrest_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/infastin/gorack/opt"
	"github.com/infastin/gorack/xrest"
)

func ExampleBindParams() {
	mux := http.NewServeMux()

	type Inline struct {
		A string `query:"a"`
		B string `query:"b"`
	}

	type Params struct {
		Inline      `inline:""`
		Foo         int              `query:"foo"`
		NullableFoo opt.NullInt[int] `query:"nullable_foo"`
		Bar         int              `header:"bar"`
		Baz         string           `path:"baz"`
	}

	mux.HandleFunc("GET /params/{baz}", func(w http.ResponseWriter, r *http.Request) {
		var params Params
		if err := xrest.BindParams(r, &params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(&params)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	uri, _ := url.Parse(ts.URL)
	query := url.Values{
		"a":   []string{"I am A"},
		"b":   []string{"Hello from B"},
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

	// Output: {"A":"I am A","B":"Hello from B","Foo":42,"NullableFoo":null,"Bar":123,"Baz":"qux"}
}
