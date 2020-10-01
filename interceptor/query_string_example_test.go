package interceptor_test

import (
	"fmt"
	"handy"
	"handy/interceptor"
	"net/http"
	"net/http/httptest"
)

type Statistics struct {
	handy.BaseHandler

	Year int `query:"year"`
}

func (h *Statistics) Get() int {
	fmt.Println(h.Year)
	return http.StatusOK
}

func NewStatistics() (handy.Handler, handy.Interceptor) {
	handler := new(Statistics)
	i := interceptor.NewIntrospector(nil, handler)
	q := interceptor.NewQueryString(i)

	return handler, q
}

func ExampleQueryString() {
	handy := handy.New()
	handy.Handle("/population", NewStatistics)
	server := httptest.NewServer(handy)
	defer server.Close()

	http.Get(server.URL + "/population?year=2020")
	// Output:
	// 2020
}
