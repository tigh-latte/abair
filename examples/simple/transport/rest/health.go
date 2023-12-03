package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/simple/domain"
)

// Health is a health handler.
type Health struct{}

// Route is a route.
func (h Health) Route(s *abair.Server) {
	s.Route("/health", func(s *abair.Server) {
		s.Use(middleware.Logger)

		s.Get("/", abair.HTTPHandlerWrapper(s, h.getHealth))
		s.Post("/", abair.HTTPHandlerWrapper(s, h.postHealth))
		s.Get("/{service}/{version}/after", abair.HTTPHandlerWrapper(s, h.serviceHealth))
		s.Get("/{arn}", abair.HTTPHandlerWrapper(s, h.arnHealth))
	})
	s.Get("/badhealth", abair.HTTPHandlerWrapper(s, h.badHealth))
}

func (h Health) getHealth(ctx context.Context, req abair.Request[struct{}, struct{}]) (struct{}, error) {
	return struct{}{}, nil
}

func (h Health) postHealth(ctx context.Context, req abair.Request[domain.HealthPostBody, struct{}]) (domain.HealthPostResponse, error) {
	return domain.HealthPostResponse{
		Service: req.Body.Service,
		Status:  "ok",
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
