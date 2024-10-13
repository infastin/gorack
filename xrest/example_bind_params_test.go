package xrest_test

import (
	"fmt"
	"net/http"

	"github.com/infastin/gorack/opt"
	"github.com/infastin/gorack/xrest"
)

func ExampleBindParams() {
	mux := http.NewServeMux()

	type Params struct {
		Foo         int              `query:"foo"`
		NullableFoo opt.NullInt[int] `query:"nullable_foo"`
		Bar         int              `header:"bar"`
		Baz         int              `path:"baz"`
	}

	mux.HandleFunc("GET /params/{baz}", func(w http.ResponseWriter, r *http.Request) {
		var params Params
		if err := xrest.BindParams(r, &params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("%#v\n", params)
	})

	http.ListenAndServe(":8080", mux)
}
