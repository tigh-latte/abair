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
		Router:       chi.NewRouter(),
		ErrorHandler: nil,
	}

	h := rest.Health{}
	h.Route(&s)

	chi.Walk(s.Router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Println(route)
		return nil
	})

	http.ListenAndServe(":3000", s)
}
