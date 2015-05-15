package interceptor

import "net/http"

type queryStringHandler interface {
	KeysWithTag(string) []string
	Field(string, string) interface{}
	Req() *http.Request
}

type QueryString struct {
	NopInterceptor

	handler queryStringHandler
}

func NewQueryString(h queryStringHandler) *QueryString {
	return &QueryString{handler: h}
}

func (q *QueryString) Before() int {
	if q.handler.Req().Form == nil {
		q.handler.Req().ParseMultipartForm(32 << 20) // 32 MB
	}

	for key, values := range q.handler.Req().Form {
		if len(values) == 0 {
			continue
		}

		if err := setValue(q.handler.Field("query", key), values[0]); err != nil {
			return http.StatusBadRequest
		}
	}

	return 0
}
