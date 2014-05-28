package handy

import "net/http"

type Interceptor interface {
	Before(w http.ResponseWriter, r *http.Request, h Handler) error
	After(w http.ResponseWriter, r *http.Request, h Handler)
	OnError(w http.ResponseWriter, r *http.Request, h Handler, err error)
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

type BeforeInterceptorFunc func(w http.ResponseWriter, r *http.Request, h Handler) error

func (i BeforeInterceptorFunc) Before(w http.ResponseWriter, r *http.Request, h Handler) error {
	return i(w, r, h)
}

func (i BeforeInterceptorFunc) After(w http.ResponseWriter, r *http.Request, h Handler)              {}
func (i BeforeInterceptorFunc) OnError(w http.ResponseWriter, r *http.Request, h Handler, err error) {}

type AfterInterceptorFunc func(w http.ResponseWriter, r *http.Request, h Handler)

func (i AfterInterceptorFunc) Before(w http.ResponseWriter, r *http.Request, h Handler) error {
	return nil
}

func (i AfterInterceptorFunc) After(w http.ResponseWriter, r *http.Request, h Handler) {
	i(w, r, h)
}

func (i AfterInterceptorFunc) OnError(w http.ResponseWriter, r *http.Request, h Handler, err error) {}

type NoErrorInterceptor struct{}

func (i NoErrorInterceptor) OnError(w http.ResponseWriter, r *http.Request, h Handler, err error) {}
