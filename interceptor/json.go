package interceptor

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type jsonHandler interface {
	Field(string, string) interface{}
	Req() *http.Request
	ResponseWriter() http.ResponseWriter
}

type JSONCodec struct {
	handler jsonHandler
}

func NewJSONCodec(h jsonHandler) *JSONCodec {
	return &JSONCodec{h}
}

func (j *JSONCodec) Before() int {
	request := j.handler.Req()
	m := strings.ToLower(request.Method)
	requestField := j.handler.Field("request", m)

	if requestField == nil {
		return 0
	}

	decoder := json.NewDecoder(request.Body)

	for {
		if err := decoder.Decode(requestField); err != nil {
			if err == io.EOF {
				break
			}

			return http.StatusInternalServerError
		}
	}

	return 0
}

func (j *JSONCodec) After(status int) int {
	headerField := j.handler.Field("response", "header")
	var header http.Header

	if headerField != nil {
		header = headerField.(http.Header)
	}

	w := j.handler.ResponseWriter()

	for k, values := range header {
		for _, value := range values {
			w.Header().Add(k, value)
		}
	}

	var response interface{}
	method := strings.ToLower(j.handler.Req().Method)

	if responseAll := j.handler.Field("response", "all"); responseAll != nil {
		response = responseAll

	} else if responseForMethod := j.handler.Field("response", method); responseForMethod != nil {
		response = responseForMethod
	}

	if response == nil {
		w.WriteHeader(status)
		return status
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)

	return status
}
