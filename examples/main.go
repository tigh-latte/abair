package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/transport/rest"
)

func main() {
	f := abair.Server{
		Router: chi.NewRouter(),
		ErrorHandler: func(_ http.ResponseWriter, _ *http.Request, err error) {
			fmt.Println(fmt.Errorf("bad error: %w", err))
		},
	}

	h := rest.Health{}
	h.Route(&f)

	http.ListenAndServe(":3000", f)
}
