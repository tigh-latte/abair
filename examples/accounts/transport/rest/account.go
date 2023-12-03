package rest

import (
	"context"
	"net/http"

	"github.com/tigh-latte/abair"
	"github.com/tigh-latte/abair/examples/accounts/domain"
)

var accounts = []domain.Account{{
	ID:   1,
	Name: "Account 1",
}, {
	ID:   2,
	Name: "Account 2",
}, {
	ID:   3,
	Name: "Account 3",
}}

type Account struct{}

func (a Account) Routes(s *abair.Server) {
	s.Route("/accounts", func(s *abair.Server) {
		s.Get("/", abair.HTTPHandlerWrapper(s, a.getAccounts))
		s.Get("/{id}", abair.HTTPHandlerWrapper(s, a.getAccount))
		s.Post("/", abair.HTTPHandlerWrapper(s, a.createAccount))
	})
}

func (a *Account) getAccounts(ctx context.Context, req abair.Request[struct{}, struct{}]) (domain.AccountGetAllResponse, error) {
	return domain.AccountGetAllResponse{Accounts: accounts}, nil
}

func (a *Account) getAccount(ctx context.Context, req abair.Request[struct{}, domain.AccountGetOnePathParams]) (domain.AccountGetOneResponse, error) {
	for _, account := range accounts {
		if account.ID == req.PathParams.ID {
			return domain.AccountGetOneResponse{Account: account}, nil
		}
	}
	return domain.AccountGetOneResponse{}, abair.NewHTTPError(
		http.StatusNotFound,
		abair.WithMessage("Account not found"),
	)
}

func (a *Account) createAccount(ctx context.Context, req abair.Request[domain.AccountCreateBody, struct{}]) (domain.AccountCreateResponse, error) {
	accounts = append(accounts, domain.Account{
		ID:   len(accounts) + 1,
		Name: req.Body.Name,
	})
	return domain.AccountCreateResponse{Account: accounts[len(accounts)-1]}, nil
}
