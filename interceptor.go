package handy

type Interceptor interface {
	Intercept(ctx *Context, h Handler)
}

type InterceptorFunc func(*Context, Handler)

func (f InterceptorFunc) Intercept(ctx *Context, h Handler) {
	f(ctx, h)
}

type InterceptorChain interface {
	Interceptors() []Interceptor
}

type InterceptorChainFunc func() []Interceptor

func (f *InterceptorFunc) Interceptors() []Interceptor {
	return f.Interceptors()
}

type NoOpInterceptorChain struct{}

func (n *NoOpInterceptorChain) Interceptors() []Interceptor {
	return nil
}
