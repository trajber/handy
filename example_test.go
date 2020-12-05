package handy_test

import (
	"fmt"
	"handy"
	"handy/interceptor"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

type Response struct {
	Kind   string
	Amount int
}

type Handler struct {
	handy.BaseHandler

	Kind     string   `urivar:"kind"`
	Amount   int      `query:"amount"`
	Response Response `response:"get"`
}

func (h *Handler) Get() int {
	h.Response.Kind = h.Kind
	h.Response.Amount = h.Amount

	return http.StatusOK
}

func NewHandler() (handy.Handler, handy.Interceptor) {
	handler := new(Handler)
	intro := interceptor.NewIntrospector(nil, handler)
	urivars := interceptor.NewURIVars(intro)
	query := interceptor.NewQueryString(urivars)
	json := interceptor.NewJSONCodec(query)

	return handler, json
}

func Example() {
	handy := handy.New()
	handy.Handle("/ball/{kind}", NewHandler)
	server := httptest.NewServer(handy)
	defer server.Close()

	response, err := http.Get(server.URL + "/ball/soccer?amount=3")

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
	// {"Kind":"soccer","Amount":3}
}
