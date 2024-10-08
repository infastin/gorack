package main

import (
	"fmt"
	"net/http"

	"github.com/infastin/gorack/fastconv"
	"github.com/infastin/gorack/opt"
	"github.com/infastin/gorack/xrest"
	"github.com/infastin/gorack/xtypes"
)

func main() {
	mux := http.NewServeMux()
	middlewares := xrest.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write(fastconv.Bytes("dummy"))
	})

	type Params struct {
		I   int              `query:"i"`
		PI  *int             `query:"pi"`
		NI  opt.NullInt[int] `query:"ni"`
		F   func()           `query:"f"`
		PF  *func()          `query:"pf"`
		AI  [2]int           `query:"ai"`
		SI  []int            `query:"si"`
		S   string           `query:"s"`
		TOD xtypes.TimeOfDay `query:"tod"`
		D   xtypes.Duration  `query:"d"`
	}

	mux.HandleFunc("GET /params", func(w http.ResponseWriter, r *http.Request) {
		var params Params
		if err := xrest.BindQuery(r, &params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("%#v\n", params)
	})

	http.ListenAndServe(":8070", middlewares(mux))
}
