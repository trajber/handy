package interceptor

import "net/http"

type QueryString struct {
	NoAfterInterceptor

	request *http.Request
	fields  map[string]interface{}
}

func NewQueryString(r *http.Request, f map[string]interface{}) *QueryString {
	return &QueryString{request: r, fields: f}
}

func (q *QueryString) Before() int {
	if q.fields == nil {
		return 0
	}

	for k, field := range q.fields {
		if field == nil {
			continue
		}

		if err := setValue(field, q.request.FormValue(k)); err != nil {
			return http.StatusBadRequest
		}
	}

	return 0
}
