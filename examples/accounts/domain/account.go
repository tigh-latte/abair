package domain

type Account struct {
	ID   int
	Name string
}

type AccountGetAllResponse struct {
	Accounts []Account
}

type AccountGetOnePathParams struct {
	ID int `path:"id"`
}

type AccountGetOneResponse struct {
	Account Account
}

type AccountCreateBody struct {
	Name string
}

type AccountCreateResponse struct {
	Account Account
}
