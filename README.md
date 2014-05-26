Handy
==========================================

Handy is a fast and simple HTTP multiplexer for Golang. It fills some gaps
related to the default Golang's HTTP multiplexer:

	* URI variable support (eg: "/test/{foo}")
	* Codecs
	* Interceptors

Handy uses the Handler As The State Of the Request. This approach allows simple and advanced usages.

## Creating a Handler
You just need to embed handy.DefaultHandler in your structure and override the HTTP method:

~~~ go
package main

import (
	"handy"
	"log"
	"net/http"
)

func main() {
	srv := handy.NewHandy()
	srv.Handle("/hello/", func() handy.Handler { return new(MyHandler) })
	log.Fatal(http.ListenAndServe(":8080", srv))
}

type MyHandler struct {
	handy.DefaultHandler
}

func (h *MyHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
~~~

# Path with variables
Path variables must be enclosed by braces.

~~~ go
srv.Handle("/hello/{name}", func() handy.Handler { 
	return new(MyHandler) 
})
~~~

And you can read them using the Handler's fields. You just need to tag the field.

~~~ go
type MyHandler struct {
	handy.DefaultHandler
	Name string `param:"name"`
}

func (h *MyHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello " + h.Name))
}
~~~

### URI variables - a complete example:
~~~ go
package main

import (
	"handy"
	"log"
	"net/http"
)

func main() {
	srv := handy.NewHandy()
	srv.Handle("/hello/{name}", func() handy.Handler {
		return new(MyHandler)
	})
	log.Fatal(http.ListenAndServe(":8080", srv))
}

type MyHandler struct {
	handy.DefaultHandler
	Name string `param:"name"`
}

func (h *MyHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello " + h.Name))
}
~~~

# Codecs
Codecs are structures that know how to unmarshal requests and marshal responses. Handy comes with a JSON codec out of the box. You just need to embed it on you Handler

~~~ go
type MyHandler struct {
	handy.DefaultHandler
	handy.JSONCodec
}
~~~

Now you're ready to create a structure that represents your protocol - in this case, using JSON tags:

~~~ go
type MyResponse struct {
	Message string `json:"message"`
}
~~~

And put it on your Handler tagged as 'codec'. This allows the JSON encoder to use this field as a response.

~~~ go
type MyHandler struct {
	handy.DefaultHandler
	handy.JSONCodec
	Response MyResponse `codec:"response"`
}
~~~

## Codecs - a complete example
~~~ go
package main

import (
	"handy"
	"log"
	"net/http"
)

func main() {
	srv := handy.NewHandy()
	srv.Handle("/hello/", func() handy.Handler {
		return new(MyHandler)
	})
	log.Fatal(http.ListenAndServe(":8080", srv))
}

type MyHandler struct {
	handy.DefaultHandler
	handy.JSONCodec
	Response MyResponse `codec:"response"`
}

func (h *MyHandler) Get(w http.ResponseWriter, r *http.Request) {
	h.Response.Message = "hello world"
}

type MyResponse struct {
	Message string `json:"message"`
}
~~~

You can create your own codecs, you just need to implement the interface:

~~~ go
type Codec interface {
	Encode(http.ResponseWriter, *http.Request, handy.Handler)
	Decode(http.ResponseWriter, *http.Request, handy.Handler)
}
~~~

If you don't need any codec you can use handy.NopCodec.

# Interceptors
To execute functions before and/or after the verb method be called you can use interceptors. To do so you need to create a InterceptorChain in you Handler to be executed Before or After the HTTP verb method.

For example: If you want to check permission

~~~ go
func (h *MyHandler) Before() handy.InterceptorChain {
	return handy.NewInterceptorChain().Chain(CheckHeader)
}

func CheckHeader(w http.ResponseWriter, r *http.Request, h handy.Handler) {
	secret := r.Header.Get("Authorization")
	if secret != "abc123" {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	}
}
~~~

## Interceptors - a complete example
~~~ go
package main

import (
	"handy"
	"log"
	"net/http"
)

func main() {
	srv := handy.NewHandy()
	srv.Handle("/hello/", func() handy.Handler {
		return new(MyHandler)
	})
	log.Fatal(http.ListenAndServe(":8080", srv))
}

type MyHandler struct {
	handy.DefaultHandler
}

func (h *MyHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Success"))
}

func (h *MyHandler) Before() handy.InterceptorChain {
	return handy.NewInterceptorChain().Chain(CheckHeader)
}

func CheckHeader(w http.ResponseWriter, r *http.Request, h handy.Handler) {
	secret := r.Header.Get("Authorization")
	if secret != "abc123" {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
	}
}
~~~

### Tests
You can use [Go's httptest package] (http://golang.org/pkg/net/http/httptest/)

~~~ go
package handler

import (
	"handy"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	mux := handy.NewHandy()
	h := new(HelloHandler)
	mux.Handle("/{name}/{id}", func() handy.Handler {
		return h
	})

	req, err := http.NewRequest("GET", "/foo/10", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if h.Id != 10 {
		t.Errorf("Unexpected Id value %d", h.Id)
	}

	if h.Name != "foo" {
		t.Errorf("Unexpected Name value %s", h.Name)
	}

	t.Logf("%d - %s - %d", w.Code, w.Body.String(), h.Id)
}
~~~
