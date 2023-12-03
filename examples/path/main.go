package main

import (
	"net/http"

	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/path/transport/rest"
)

func main() {
	s := abair.NewServer()

	s.Route("/api/v1", func(s *abair.Server) {
		(&rest.Schema{}).Routes(s)
	})

	http.ListenAndServe(":8080", s)
}
