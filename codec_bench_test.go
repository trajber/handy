package handy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestStruct struct {
	Name string `json:"name" method:"get"`
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

	codec := JSONCodec{}

	handler := new(struct {
		DefaultHandler
		Request TestStruct `codec:"request" method:"get"`
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codec.Decode(w, req, handler)
		if handler.Request.Id != 10 {
			b.Fail()
		}
	}
}
