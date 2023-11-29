package rest

import (
	"context"

	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/domain"
)

type Health struct{}

func (h Health) Route(s *abair.Server) {
	abair.Get(s, "/health", h.getHealth)
	abair.Post(s, "/health", h.postHealth)
}

func (h Health) getHealth(ctx context.Context, req abair.Request[struct{}]) (struct{}, error) {
	return struct{}{}, nil
}

func (h Health) postHealth(ctx context.Context, req abair.Request[domain.HealthPostRequest]) (domain.HealthPostResponse, error) {
	return domain.HealthPostResponse{
		Status: "ok",
	}, nil
}
