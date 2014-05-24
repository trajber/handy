package handy

type Interceptor interface {
	Intercept(ctx *Context, h Handler)
}

type InterceptorFunc func(*Context, Handler)

func (f InterceptorFunc) Intercept(ctx *Context, h Handler) {
	f(ctx, h)
}
