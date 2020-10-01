package interceptor

import (
	"handy"
	"net/http"
)

// URIVars automatically converts Handy's URI variables into fields of the
// provided struct. The struct is set via Introspector and thus, that
// interceptor must be run before this one.
//
// The struct's fields must be tagged in the format `urivars:"name"`, where
// name is the name of the URI variable defined when registering the handler.
type URIVars interface {
	Introspector
}

// URIVarsAPI is the API provided by URIVars to be used by other interceptors.
type URIVarsAPI interface {
	IntrospectorAPI
}

type uriVars struct {
	handy.BaseInterceptor
	IntrospectorAPI
}

// NewURIVars creates a URIVars. It uses the API provided by Introspector, and
// thus requires as argument any interceptor compatible with the Introspector
// interface.
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
