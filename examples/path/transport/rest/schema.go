package rest

import (
	"context"

	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/path/domain"
)

type Schema struct{}

func (s Schema) Routes(server *abair.Server) {
	abair.Route(server, "/schema", func(server *abair.Server) {
		abair.Get(server, "/1/{person}", s.handle1)
		abair.Get(server, "/2/{person}", s.handle2)
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
