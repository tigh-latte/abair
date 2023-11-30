package domain

type HealthGetResponse struct {
	Status string
}

type HealthGetServicePath struct {
	Service string `path:"service"`
	Version int    `path:"version"`
}

type HealthGetServiceResponse struct {
	Service string
	Version int
}

type HealthPostRequest struct {
	Service string
}

type HealthPostResponse struct {
	Status string
}
