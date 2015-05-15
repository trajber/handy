package interceptor

import (
	"net/http"

	"github.com/gustavo-hms/handy"
)

type uriVarsHandler interface {
	URIVars() handy.URIVars
	Field(string, string) interface{}
}

type URIVars struct {
	NoAfterInterceptor

	handler uriVarsHandler
}

func NewURIVars(h uriVarsHandler) *URIVars {
	return &URIVars{handler: h}
}

func (u *URIVars) Before() int {
	for k, value := range u.handler.URIVars() {
		field := u.handler.Field("urivar", k)

		if field == nil {
			continue
		}

		if err := setValue(field, value); err != nil {
			return http.StatusBadRequest
		}
	}

	return 0
}
