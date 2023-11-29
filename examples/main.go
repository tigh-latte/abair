package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/transport/rest"
)

func main() {
	s := abair.Server{
		Router: chi.NewRouter(),
		ErrorHandler: func(_ http.ResponseWriter, _ *http.Request, err error) {
			fmt.Println(fmt.Errorf("bad error: %w", err))
		},
	}

	h := rest.Health{}
	h.Route(&s)

	http.ListenAndServe(":3000", s)
}
