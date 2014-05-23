package handy

import (
	"net/http"
)

type Handler interface {
	Get(ctx *Context) (int, error)
	Post(ctx *Context) (int, error)
	Put(ctx *Context) (int, error)
	Delete(ctx *Context) (int, error)
	Patch(ctx *Context) (int, error)
}

type DefaultHandler struct {
	Handler http.Handler
}

func (s *DefaultHandler) defaultHandler(ctx *Context) (int, error) {
	if s.Handler != nil {
		s.Handler.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	} else {
		ctx.ResponseWriter.WriteHeader(http.StatusNotImplemented)
	}

	return http.StatusNotImplemented, nil
}

func (s *DefaultHandler) Get(ctx *Context) (int, error) {
	return s.defaultHandler(ctx)
}

func (s *DefaultHandler) Post(ctx *Context) (int, error) {
	return s.defaultHandler(ctx)
}

func (s *DefaultHandler) Put(ctx *Context) (int, error) {
	return s.defaultHandler(ctx)
}

func (s *DefaultHandler) Delete(ctx *Context) (int, error) {
	return s.defaultHandler(ctx)
}

func (s *DefaultHandler) Patch(ctx *Context) (int, error) {
	return s.defaultHandler(ctx)
}
