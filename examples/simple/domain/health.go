package domain

import (
	"fmt"
	"strings"
)

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

// HealthGetArnPath is a health get arn path.
type HealthGetArnPath struct {
	Arn ARN `path:"arn"`
}

// HealthGetArnResponse is a health get arn response.
type HealthGetArnResponse struct {
	Partition string
	Region    string
	AccountID string
}

// ARN is an arn.
type ARN struct {
	Partition string
	Region    string
	AccountID string
}

func (a *ARN) ParsePath(s string) error {
	if count := strings.Count(s, ":"); count != 2 {
		return fmt.Errorf("invalid ARN: %q %d", s, count)
	}
	fn := func(oo ...*string) {
		target := s

		for i := range oo {
			o := oo[i]
			idx := strings.Index(target, ":")
			if idx == -1 {
				*o = target
				return
			}

			*o = target[:idx]
			target = target[idx+1:]
		}
	}

	fn(&a.Partition, &a.Region, &a.AccountID)

	return nil
}
