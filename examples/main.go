package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/transport/rest"
)

func main() {
	s := abair.NewServer()

	abair.Route(s, "/api/v1", func(s *abair.Server) {
		abair.Use(s,
			middleware.RequestID,
			middlewareExample(s.Logger),
		)

		h := rest.Health{}
		h.Route(s)
	})

	http.ListenAndServe(":3000", s)
}

func middlewareExample(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.LogAttrs(r.Context(), slog.LevelInfo, "hi")
			next.ServeHTTP(w, r)
			log.LogAttrs(r.Context(), slog.LevelInfo, "bye")
		})
	}
}
