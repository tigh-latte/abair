package main

import (
	"net/http"

	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/accounts/transport/rest"
)

func main() {
	server := abair.NewServer()

	abair.Route(server, "/api/v1", func(s *abair.Server) {
		(&rest.Account{}).Routes(s)
	})

	http.ListenAndServe(":8080", server)
}
