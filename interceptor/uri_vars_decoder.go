package interceptor

import (
	"handy"
	"net/http"
)

type uriVarsHandler interface {
	URIVars() handy.URIVars
	Field(string, string) interface{}
}

type URIVars interface {
	Introspector
}

type URIVarsAPI interface {
	IntrospectorAPI
}

type uriVars struct {
	handy.ProtoInterceptor
	IntrospectorAPI
}

func NewURIVars(previous Introspector) URIVars {
	if previous == nil {
		panic("URIVars' dependency can not be nil")
	}

	u := &uriVars{IntrospectorAPI: previous}
	u.SetPrevious(previous)

	return u
}

func (u *uriVars) Before() int {
	for k, value := range u.URIVars {
		field := u.Field("urivar", k)

		if field == nil {
			continue
		}

		if err := setValue(field, value); err != nil {
			return http.StatusBadRequest
		}
	}

	return 0
}
