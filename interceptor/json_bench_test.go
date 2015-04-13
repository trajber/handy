package interceptor

import (
	"handy"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestStruct struct {
	Name string `json:"name"`
	Id   int    `json:id`
}

var (
	payload = `
{
	"name":"foo",
	"id":10
}`
)

func BenchmarkDecodeJSON(b *testing.B) {
	req, err := http.NewRequest("GET", "/", strings.NewReader(payload))
	if err != nil {
		b.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := new(struct {
		IntrospectorEmbedded
		handy.DefaultHandler
		Request TestStruct `request:"get"`
	})
	handler.DefaultHandler = handy.BuildDefaultHandler(w, req, nil)
	codec := NewJSONCodec(handler, w, req)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codec.Before()
		if handler.Request.Id != 10 {
			b.Fail()
		}
	}
}
