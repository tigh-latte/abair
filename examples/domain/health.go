package domain

// HealthGetResponse is a health get response.
type HealthGetResponse struct {
	Status string
}

// HealthGetServicePath is a health get service path.
type HealthGetServicePath struct {
	Service string `path:"service"`
	Version int    `path:"version"`
}

// HealthGetServiceResponse is a health get service response.
type HealthGetServiceResponse struct {
	Service string
	Version int
}

// HealthPostBody is a health post request.
type HealthPostBody struct {
	Service string
}

// HealthPostResponse is a health post response.
type HealthPostResponse struct {
	Status string
}
