package domain

type HealthGetResponse struct {
	Status string
}

type HealthPostRequest struct {
	Service string
}

type HealthPostResponse struct {
	Status string
}
