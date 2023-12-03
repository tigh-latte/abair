package rest

import (
	"context"

	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/path/domain"
)

type Schema struct{}

func (s Schema) Routes(svr *abair.Server) {
	svr.Route("/schema", func(svr *abair.Server) {
		svr.Get("/1/{person}", abair.HTTPHandlerWrapper(svr, s.handle1))
		svr.Get("/2/{person}", abair.HTTPHandlerWrapper(svr, s.handle2))
	})
}

func (s Schema) handle1(ctx context.Context, req abair.Request[struct{}, domain.Handle1Path]) (domain.Handle1Response, error) {
	return domain.Handle1Response{
		Name: req.PathParams.Person.Name,
		Age:  req.PathParams.Person.Age,
	}, nil
}

func (s Schema) handle2(ctx context.Context, req abair.Request[struct{}, domain.Handle2Path]) (domain.Handle2Response, error) {
	return domain.Handle2Response{
		Name: req.PathParams.Person.Name,
		Age:  req.PathParams.Person.Age,
	}, nil
}
