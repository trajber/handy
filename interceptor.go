package handy

type Interceptor interface {
	Before() int
	After(int) int
}

type InterceptorChain []Interceptor

func (c InterceptorChain) Chain(f Interceptor) InterceptorChain {
	return append(c, f)
}

func NewInterceptorChain() InterceptorChain {
	return make([]Interceptor, 0)
}

type NopInterceptorChain struct{}

func (n *NopInterceptorChain) Interceptors() InterceptorChain {
	return NewInterceptorChain()
}
