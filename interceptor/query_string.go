package interceptor

import (
	"handy"
	"net/http"
)

type QueryString interface {
	Introspector
}

type QueryStringAPI interface {
	IntrospectorAPI
}

type queryString struct {
	handy.BaseInterceptor
	IntrospectorAPI
}

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
