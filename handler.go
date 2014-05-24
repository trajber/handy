package handy

import "net/http"

type Handler interface {
	Get(ctx *Context)
	Post(ctx *Context)
	Put(ctx *Context)
	Delete(ctx *Context)
	Patch(ctx *Context)
	Decode(*Context, Handler)
	Encode(*Context, Handler)
	Interceptors() []Interceptor
}

type DefaultHandler struct {
	http.Handler
	NoOpCodec
	NoOpInterceptorChain
}

func (s *DefaultHandler) defaultHandler(ctx *Context) {
	if s.Handler != nil {
		s.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	} else {
		ctx.ResponseWriter.WriteHeader(http.StatusNotImplemented)
	}
}

func (s *DefaultHandler) Get(ctx *Context) {
	s.defaultHandler(ctx)
}

func (s *DefaultHandler) Post(ctx *Context) {
	s.defaultHandler(ctx)
}

func (s *DefaultHandler) Put(ctx *Context) {
	s.defaultHandler(ctx)
}

func (s *DefaultHandler) Delete(ctx *Context) {
	s.defaultHandler(ctx)
}

func (s *DefaultHandler) Patch(ctx *Context) {
	s.defaultHandler(ctx)
}

type JSONHandler struct {
	DefaultHandler
	JSONCodec
}
