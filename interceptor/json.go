package interceptor

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type jsonHandler interface {
	Field(string, string) interface{}
}

type JSONCodec struct {
	handler  jsonHandler
	response http.ResponseWriter
	request  *http.Request
}

func NewJSONCodec(h jsonHandler, w http.ResponseWriter, r *http.Request) *JSONCodec {
	return &JSONCodec{handler: h, response: w, request: r}
}

func (j *JSONCodec) Before() int {
	m := strings.ToLower(j.request.Method)
	requestField := j.handler.Field("request", m)

	if requestField == nil {
		return 0
	}

	decoder := json.NewDecoder(j.request.Body)

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

	for k, values := range header {
		for _, value := range values {
			j.response.Header().Add(k, value)
		}
	}

	var response interface{}
	method := strings.ToLower(j.request.Method)

	if responseAll := j.handler.Field("response", "all"); responseAll != nil {
		response = responseAll

	} else if responseForMethod := j.handler.Field("response", method); responseForMethod != nil {
		response = responseForMethod
	}

	if response == nil {
		j.response.WriteHeader(status)
		return status
	}

	j.response.Header().Set("Content-Type", "application/json")
	j.response.WriteHeader(status)
	json.NewEncoder(j.response).Encode(response)

	return status
}
