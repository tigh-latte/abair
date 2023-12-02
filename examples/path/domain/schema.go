package domain

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/tigh-latte/abair"
)

type Person struct {
	Name string
	Age  int64
}

type (
	Handle1Path struct {
		Person Person `path:"person"`
	}
	Handle1Response struct {
		Name string
		Age  int64
	}
)

type (
	Handle2Path struct {
		Person *Person `path:"person"`
	}
	Handle2Response struct {
		Name string
		Age  int64
	}
)

func (p *Person) ParsePath(s string) error {
	split := strings.LastIndex(s, ":")
	p.Name = s[0:split]
	age, err := strconv.ParseInt(s[split+1:], 10, 32)
	if err != nil {
		return abair.NewHTTPError(http.StatusBadRequest,
			abair.WithMessage("invalid age"),
			abair.WithInternal(err),
		)
	}
	p.Age = age
	return nil
}
