package handy

import "net/http"

type Interceptor interface {
	Before(w http.ResponseWriter, r *http.Request)
	After(w http.ResponseWriter, r *http.Request)
}

type InterceptorChain []Interceptor

func (c InterceptorChain) Chain(f Interceptor) InterceptorChain {
	c = append(c, f)
	return c
}

func NewInterceptorChain() InterceptorChain {
	return make([]Interceptor, 0)
}

type NopInterceptorChain struct{}

func (n *NopInterceptorChain) Interceptors() InterceptorChain {
	return NewInterceptorChain()
}

type BeforeInterceptorFunc func(w http.ResponseWriter, r *http.Request)

func (i BeforeInterceptorFunc) Before(w http.ResponseWriter, r *http.Request) {
	i(w, r)
}

func (i BeforeInterceptorFunc) After(w http.ResponseWriter, r *http.Request) {}

type AfterInterceptorFunc func(w http.ResponseWriter, r *http.Request)

func (i AfterInterceptorFunc) Before(w http.ResponseWriter, r *http.Request) {}

func (i AfterInterceptorFunc) After(w http.ResponseWriter, r *http.Request) {
	i(w, r)
}
