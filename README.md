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
	"github.com/trajber/handy"
	"log"
	"net/http"
)

func main() {
	srv := handy.NewHandy()
	srv.Handle("/hello", func() handy.Handler {
		return &MyHandler{}
	})
	log.Fatal(http.ListenAndServe(":8080", srv))
}

type MyHandler struct {
	handy.DefaultHandler
}

func (h *MyHandler) Get() int {
	h.ResponseWriter().Write([]byte("Hello World"))
	return http.StatusOK
}
~~~

# Interceptors
The true power of this framework comes from the use of interceptors. They are special units that are called before and after every handler method call. With interceptors, one can automate most of the repetitive tasks involving a request handling, like the setup and commit of a database transaction, JSON serialisation and automatic decode of URI parameters.

If the handler register the interceptor chain `[a, b, c]`, the framework will call, in order:
~~~
a.Before
b.Before
c.Before
handler.Method (any of Get, Put, Post, Delete or Patch)
c.After
b.After
a.After
~~~
If any of the interceptors' Before returns a code different than zero, the chain is interrupted. Say, for example, that the Before method of the b interceptor returns http.StatusInternalServerError; then, the execution will be:
~~~
a.Before
b.Before
b.After
a.After
~~~

## Interceptors - simple example
~~~ go
package main

import (
	"github.com/trajber/handy"
	"log"
	"net/http"
	"time"
)

func main() {
	srv := handy.NewHandy()
	srv.Handle("/hello", func() handy.Handler {
		return &MyHandler{}
	})
	log.Fatal(http.ListenAndServe(":8080", srv))
}

type MyHandler struct {
	handy.DefaultHandler
}

func (h *MyHandler) Get() int {
	h.ResponseWriter().Write([]byte("Hello World"))
	return http.StatusOK
}

func (h *MyHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().Chain(new(TimerInterceptor))
}

type TimerInterceptor struct {
	Timer time.Time
}

func (i *TimerInterceptor) Before() int {
	i.Timer = time.Now()
	return 0
}

func (i *TimerInterceptor) After(status int) int {
	log.Println("Took", time.Since(i.Timer))
	return status
}
~~~

## JSON Codec interceptor
Handy comes with a JSONCodec interceptor out of the box. It can be used to automatically unmarshal requests and marshal responses using JSON. It does so by reading special tags in your handler:

~~~ go
type MyResponse struct {
	Message string `json:"message"`
}

type MyHandler struct {
	handy.DefaultHandler
	interceptor.IntrospectorCompliant

	// this structure will be used only for GET and PUT methods
	Response MyResponse `response:"get,put"`
}
~~~

Now, you just need to include JSONCodec in the handler's interceptor chain:
~~~ go
func (h *MyHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(interceptor.NewIntrospector(h)).
		Chain(interceptor.NewJSONCodec(h))
}
~~~

### JSON Codec interceptor - a complete example:
~~~ go
package main

import (
	"github.com/trajber/handy"
	"github.com/trajber/handy/interceptor"
    "log"
    "net/http"
)

func main() {
    srv := handy.NewHandy()
	srv.Handle("/hello", func() handy.Handler {
		return &MyHandler{}
	})
    log.Fatal(http.ListenAndServe(":8080", srv))
}

type MyHandler struct {
    handy.DefaultHandler
	interceptor.IntrospectorCompliant

    Response MyResponse `response:"all"`
}

func (h *MyHandler) Get() int {
    h.Response.Message = "hello world"
	return http.StatusOK
}

func (h *MyHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(interceptor.NewIntrospector(h)).
		Chain(interceptor.NewJSONCodec(h))
}

type MyResponse struct {
    Message string `json:"message"`
}
~~~

### JSON Codec interceptor - An example with JSON in request and response:
~~~go
package main

import (
	"log"
	"net/http"

	"github.com/trajber/handy"
	"github.com/trajber/handy/interceptor"
)

func main() {
	srv := handy.NewHandy()
	srv.Handle("/hello", func() handy.Handler {
		return &MyHandler{}
	})
	log.Fatal(http.ListenAndServe(":8181", srv))
}

type MyHandler struct {
	handy.DefaultHandler
	interceptor.IntrospectorCompliant

	Response MyResponse `response:"post"`
	Request  MyRequest  `request:"post"`
}

func (h *MyHandler) Post() int {
	h.Response.Answer = "You asked me about " + h.Request.Question
	return http.StatusOK
}

func (h *MyHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(interceptor.NewIntrospector(h)).
		Chain(interceptor.NewJSONCodec(h))
}

type MyResponse struct {
	Answer string `json:"answer"`
}

type MyRequest struct {
	Question string `json:"question"`
}
~~~
When you submit the content using the HTTP verb POST:
~~~javascript
{
	"question": "life"
}
~~~
You will get:
~~~javascript
{
	"answer": "You asked me about life"
}
~~~

## URIVar interceptor
Handy can automatically set the URI parameters in the handler using the included URIVar interceptor. It has support for Go native types plus any type that implements the TextUnmarshaler interface:

~~~go
package main

import (
	"log"
	"net"
	"net/http"

	"github.com/trajber/handy"
	"github.com/trajber/handy/interceptor"
)

func main() {
srv := handy.NewHandy()
	srv.Handle("/user/{user}/machine/{ip}", func() handy.Handler {
		return &MyHandler{}
	})
	log.Fatal(http.ListenAndServe(":8181", srv))
}

type MyHandler struct {
	handy.DefaultHandler
	interceptor.IntrospectorCompliant

	User     string     `urivar:"user"`
	IP       net.IP     `urivar:"ip"`
	Response MyResponse `response:"get"`
}

func (h *MyHandler) Get() int {
	h.Response.Message = "Request from user " + h.User + " at " + h.IP.String()
	return http.StatusOK
}

func (h *MyHandler) Interceptors() handy.InterceptorChain {
	return handy.NewInterceptorChain().
		Chain(interceptor.NewIntrospector(h)).
		Chain(interceptor.NewJSONCodec(h)).
		Chain(interceptor.NewURIVars(h))
}

type MyResponse struct {
	Message string `json:"message"`
}
~~~

Thanks to the fact that the URIVar interceptor supports TextUnmarshalers, one can use it to automatically validate a URI variable before it reaches the handler:
~~~go
type limitedString string

func (l *limitedString) UnmarshalText(data []byte) error {
	text := string(data)
	length := len(text)

	if length < 5 || length > 50 {
		return errors.New("Wrong size for parameter")
	}

	*l = text
	return nil
}
~~~

You can do the same with the QueryString interceptor, also included in the handy/interceptor package.

#Logging
Bad things happens even inside Handy; You can set your own function to handle Handy errors.

~~~go
package main

import (
	"github.com/trajber/handy"
	"log"
	"net/http"
)

func main() {
    srv := handy.NewHandy()
    // This function will be called when
    // some error occurs inside Handy code.
    srv.ErrorFunc = func(e error) {
    	// here you can handle the error
    	log.Println(e)
  	}
	srv.Handle("/hello", func() handy.Handler {
		return &MyHandler{}
	})

    log.Fatal(http.ListenAndServe(":8080", srv))
}
~~~
