package rest

import (
	"context"

	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/domain"
)

type Health struct{}

func (h Health) Route(f *abair.Server) {
	abair.Get(f, "/health", h.getHealth)
	abair.Post(f, "/health", h.postHealth)
}

func (h Health) getHealth(ctx context.Context, req abair.Request[struct{}]) (struct{}, error) {
	return struct{}{}, nil
}

func (h Health) postHealth(ctx context.Context, req abair.Request[domain.HealthPostRequest]) (domain.HealthPostResponse, error) {
	return domain.HealthPostResponse{
		Status: "ok",
	}, nil
}
