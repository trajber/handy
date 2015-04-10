package interceptor

import (
	"net/http"
	"reflect"
)

type queryStringHandler interface {
	FieldValues(tag string) (map[string]reflect.Value, bool)
	Field(string, string) interface{}
	Req() *http.Request
}

type QueryString struct {
	NoAfterInterceptor

	handler queryStringHandler
}

func NewQueryString(h queryStringHandler) *QueryString {
	return &QueryString{handler: h}
}

func (q *QueryString) Before() int {
	queries, ok := q.handler.FieldValues("query")

	if !ok {
		return 0
	}

	request := q.handler.Req()

	for k := range queries {
		field := q.handler.Field("query", k)

		if field == nil {
			continue
		}

		if err := setValue(field, request.FormValue(k)); err != nil {
			return http.StatusBadRequest
		}
	}

	return 0
}
