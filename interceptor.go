package handy

import "net/http"

type Interceptor func(w http.ResponseWriter, r *http.Request, h Handler)

type InterceptorChain []Interceptor

func (c InterceptorChain) Chain(f Interceptor) InterceptorChain {
	c = append(c, f)
	return c
}

func NewInterceptorChain() InterceptorChain {
	return make([]Interceptor, 0)
}

type NopInterceptorChain struct{}

func (n *NopInterceptorChain) Before() InterceptorChain {
	return NewInterceptorChain()
}

func (n *NopInterceptorChain) After() InterceptorChain {
	return NewInterceptorChain()
}
