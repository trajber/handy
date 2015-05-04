package interceptor

import (
	"net/http"

	"github.com/gustavo-hms/handy"
)

type URIVars struct {
	NoAfterInterceptor

	uriVars handy.URIVars
	fields  map[string]interface{}
}

func NewURIVars(u handy.URIVars, f map[string]interface{}) *URIVars {
	return &URIVars{uriVars: u, fields: f}
}

func (u *URIVars) Before() int {
	for k, value := range u.uriVars {
		field := u.fields[k]

		if field == nil {
			continue
		}

		if err := setValue(field, value); err != nil {
			return http.StatusBadRequest
		}
	}

	return 0
}
