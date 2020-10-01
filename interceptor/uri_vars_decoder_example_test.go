package interceptor_test

import (
	"fmt"
	"handy"
	"handy/interceptor"
	"net/http"
	"net/http/httptest"
)

type Handler struct {
	handy.BaseHandler

	Year int `urivar:"year"`
}

func (h *Handler) Get() int {
	fmt.Println(h.Year)
	return http.StatusOK
}

func NewHandler() (handy.Handler, handy.Interceptor) {
	handler := new(Handler)
	i := interceptor.NewIntrospector(nil, handler)
	u := interceptor.NewURIVars(i)

	return handler, u
}

func ExampleURIVars() {
	handy := handy.New()
	handy.Handle("/population/{year}", NewHandler)
	server := httptest.NewServer(handy)
	defer server.Close()

	http.Get(server.URL + "/population/2020")
	// Output:
	// 2020
}
