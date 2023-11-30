package rest

import (
	"context"
	"net/http"

	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/domain"
)

type Health struct{}

func (h Health) Route(s *abair.Server) {
	abair.Route(s, "/health", func(s *abair.Server) {
		abair.Get(s, "/", h.getHealth)
		abair.Post(s, "/", h.postHealth)
		abair.Get(s, "/{service}", h.serviceHealth)
	})
	abair.Get(s, "/badhealth", h.badHealth)
}

func (h Health) getHealth(ctx context.Context, req abair.Request[struct{}]) (struct{}, error) {
	return struct{}{}, nil
}

func (h Health) postHealth(ctx context.Context, req abair.Request[domain.HealthPostRequest]) (domain.HealthPostResponse, error) {
	return domain.HealthPostResponse{
		Status: "ok",
	}, nil
}

func (h Health) badHealth(ctx context.Context, req abair.Request[struct{}]) (struct{}, error) {
	return struct{}{}, abair.NewHTTPError(http.StatusBadRequest, abair.WithMessage("bad health"))
}

func (h Health) serviceHealth(ctx context.Context, req abair.Request[struct{}]) (domain.HealthGetServiceResponse, error) {
	return domain.HealthGetServiceResponse{
		Service: req.PathParams["service"],
	}, nil
}
