package handy

import "net/http"

type Interceptor interface {
	Intercept(w http.ResponseWriter, r *http.Request, h Handler)
}

type InterceptorFunc func(w http.ResponseWriter, r *http.Request, h Handler)

func (f InterceptorFunc) Intercept(w http.ResponseWriter,
	r *http.Request,
	h Handler) {
	f(w, r, h)
}

type InterceptorChain interface {
	Interceptors() []Interceptor
}

type InterceptorChainFunc func() []Interceptor

func (f *InterceptorFunc) Interceptors() []Interceptor {
	return f.Interceptors()
}

type NoOpInterceptorChain struct{}

func (n *NoOpInterceptorChain) Before() []Interceptor {
	return nil
}

func (n *NoOpInterceptorChain) After() []Interceptor {
	return nil
}
