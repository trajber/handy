package handy_test

import (
	"handy"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testInterceptorHandler struct {
	handy.ProtoHandler
}

func NewTestInterceptorHandler() (handy.Handler, handy.Interceptor) {
	return new(testInterceptorHandler), new(dummyinterceptor)
}

type dummyinterceptor struct {
	handy.ProtoInterceptor
}

func (i *dummyinterceptor) Before() int {
	for j := 0; j < 10000; j++ {
	}

	return 0
}

func (i *dummyinterceptor) After(int) int {
	for j := 0; j < 10000; j++ {
	}

	return 0
}

func BenchmarkInterceptorExecution(b *testing.B) {
	mux := handy.New()
	mux.Handle("/foo", NewTestInterceptorHandler)

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
