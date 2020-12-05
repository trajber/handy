package interceptor_test

import (
	"fmt"
	"handy"
	"handy/interceptor"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

type Sum struct {
	N int
	M int
}

type Result struct {
	Sum int
}

type Math struct {
	handy.BaseHandler

	Request  Sum    `request:"get,post"`
	Response Result `response:"get,post"`
}

func (h *Math) Post() int {
	h.Response.Sum = h.Request.N + h.Request.M
	return http.StatusOK
}

func NewMath() (handy.Handler, handy.Interceptor) {
	handler := new(Math)
	intro := interceptor.NewIntrospector(nil, handler)
	json := interceptor.NewJSONCodec(intro)

	return handler, json
}

func ExampleJSONCodec() {
	handy := handy.New()
	handy.Handle("/math/sum", NewMath)
	server := httptest.NewServer(handy)
	defer server.Close()

	response, err := http.Post(
		server.URL+"/math/sum",
		"application/json",
		strings.NewReader(`{"N":17,"M":19}`),
	)

	if err != nil {
		log.Fatal(err)
	}

	result, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", result)
	// Output:
	// {"Sum":36}
}
