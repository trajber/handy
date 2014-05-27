package handy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestHandler struct {
	DefaultHandler
}

func BenchmarkSimpleRequest(b *testing.B) {
	mux := NewHandy()
	mux.Handle("/foo", func() Handler { return new(TestHandler) })

	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		b.Fatal(err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkPathWithVariable(b *testing.B) {
	mux := NewHandy()
	mux.Handle("/foo/{name}", func() Handler { return new(TestHandler) })

	req, err := http.NewRequest("GET", "/foo/bar", nil)
	if err != nil {
		b.Fatal(err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkPathWithVariables(b *testing.B) {
	mux := NewHandy()
	mux.Handle("/foo/{name}/{age}/{nono}", func() Handler { return new(TestHandler) })

	req, err := http.NewRequest("GET", "/foo/bar/100/x", nil)
	if err != nil {
		b.Fatal(err)
	}

	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(w, req)
	}
}
