package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/domain"
)

// Health is a health handler.
type Health struct{}

// Route is a route.
func (h Health) Route(s *abair.Server) {
	abair.Route(s, "/health", func(s *abair.Server) {
		abair.Use(s, func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				s.Logger.LogAttrs(r.Context(), slog.LevelInfo, "inline hi")
				next.ServeHTTP(w, r)
				s.Logger.LogAttrs(r.Context(), slog.LevelInfo, "inline bye")
			})
		})
		abair.Get(s, "/", h.getHealth)
		abair.Post(s, "/", h.postHealth)
		abair.Get(s, "/{service}/{version}/after", h.serviceHealth)
	})
	abair.Get(s, "/badhealth", h.badHealth)
}

func (h Health) getHealth(ctx context.Context, req abair.Request[struct{}, struct{}]) (struct{}, error) {
	return struct{}{}, nil
}

func (h Health) postHealth(ctx context.Context, req abair.Request[domain.HealthPostBody, struct{}]) (domain.HealthPostResponse, error) {
	return domain.HealthPostResponse{
		Status: "ok",
	}, nil
}

func (h Health) badHealth(ctx context.Context, req abair.Request[struct{}, struct{}]) (struct{}, error) {
	return struct{}{}, abair.NewHTTPError(http.StatusBadRequest, abair.WithMessage("bad health"))
}

func (h Health) serviceHealth(ctx context.Context, req abair.Request[struct{}, domain.HealthGetServicePath]) (domain.HealthGetServiceResponse, error) {
	return domain.HealthGetServiceResponse{
		Service: req.PathParams.Service,
		Version: req.PathParams.Version,
	}, nil
}
