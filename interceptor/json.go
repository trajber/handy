package interceptor

import (
	"encoding/json"
	"handy"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// JSONCodec automatically handles marshaling and unmarshaling data. It uses
// the API provided by Introspector to find fields tagged with a
// `request:"method"` or a `response:"method"` tag, where method can be any
// combination of get, put, post, patch, delete separated by commas. It will
// then unmarshal a request's body into the corresponding field or marshal the
// value found in the field into the response's body.
type JSONCodec interface {
	Introspector
}

// JSONCodecAPI is the API provided by JSONCodec to be used by other interceptors.
type JSONCodecAPI interface {
	IntrospectorAPI
}

type jsonCodec struct {
	handy.BaseInterceptor
	IntrospectorAPI
}

// NewJSONCodec creates a JSONCodec. It uses the API provided by Introspector,
// and thus requires as argument any interceptor compatible with the
// Introspector interface.
func NewJSONCodec(previous Introspector) JSONCodec {
	if previous == nil {
		panic("JSONCodec's dependency can not be nil")
	}

	j := &jsonCodec{IntrospectorAPI: previous}
	j.SetPrevious(previous)

	return j
}

func (j *jsonCodec) Before() int {
	m := strings.ToLower(j.Request.Method)
	requestField := j.Field("request", m)

	if requestField == nil {
		return 0
	}

	decoder := json.NewDecoder(j.Request.Body)

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

func (j *jsonCodec) After(status int) int {
	headerField := j.Field("response", "header")

	if headerField != nil {
		if header, ok := headerField.(*http.Header); ok {
			for k, values := range *header {
				for _, value := range values {
					j.ResponseWriter.Header().Add(k, value)
				}
			}
		}
	}

	var response interface{}
	method := strings.ToLower(j.Request.Method)

	if responseAll := j.Field("response", "all"); responseAll != nil {
		response = responseAll

	} else if responseForMethod := j.Field("response", method); responseForMethod != nil {
		response = responseForMethod
	}

	var buf []byte
	buf, err := json.Marshal(response)
	if err != nil || response == nil {
		j.ResponseWriter.WriteHeader(status)
		return status
	}

	j.ResponseWriter.Header().Set("Content-Type", "application/json")
	j.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	j.ResponseWriter.WriteHeader(status)
	j.ResponseWriter.Write(buf)

	return status
}
