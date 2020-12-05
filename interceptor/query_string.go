package interceptor

import (
	"handy"
	"net/http"
)

// QueryString sets fields of a struct from values of query string arguments.
// The fields must be tagged in the following format: `query:"name"`, where
// name is a query string parameter. It uses the API provided by Introspector
// to find such fields, and thus must be run after that one.
type QueryString interface {
	Introspector
}

// QueryStringAPI is the API provided by QueryString to be used by other
// interceptors.
type QueryStringAPI interface {
	IntrospectorAPI
}

type queryString struct {
	handy.BaseInterceptor
	IntrospectorAPI
}

// NewQueryString creates a QueryString. It uses the API provided by
// Introspector, and thus requires as argument any interceptor compatible with
// the Introspector interface.
func NewQueryString(previous Introspector) QueryString {
	if previous == nil {
		panic("QueryString's dependency can not be nil")
	}

	q := &queryString{IntrospectorAPI: previous}
	q.SetPrevious(previous)

	return q
}

func (q *queryString) Before() int {
	if q.Request.Form == nil {
		q.Request.ParseMultipartForm(32 << 20) // 32 MB
	}

	for key, values := range q.Request.Form {
		if len(values) == 0 {
			continue
		}

		if err := setValue(q.Field("query", key), values[0]); err != nil {
			return http.StatusBadRequest
		}
	}

	return 0
}
