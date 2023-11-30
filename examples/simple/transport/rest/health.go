package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/domain"
)

// Health is a health handler.
type Health struct{}

// Route is a route.
func (h Health) Route(s *abair.Server) {
	abair.Route(s, "/health", func(s *abair.Server) {
		abair.Use(s,
			middleware.Logger,
		)

		abair.Get(s, "/", h.getHealth)
		abair.Post(s, "/", h.postHealth)
		abair.Get(s, "/{service}/{version}/after", h.serviceHealth)
		abair.Get(s, "/{arn}", h.arnHealth)
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

func (h Health) arnHealth(ctx context.Context, req abair.Request[struct{}, domain.HealthGetArnPath]) (domain.HealthGetArnResponse, error) {
	return domain.HealthGetArnResponse{
		Partition: req.PathParams.Arn.Partition,
		Region:    req.PathParams.Arn.Region,
		AccountID: req.PathParams.Arn.AccountID,
	}, nil
}
