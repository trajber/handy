package mux

import (
	"net/http"
)

type Service interface {
	Get(ctx *Context)
	Post(ctx *Context)
	Put(ctx *Context)
	Delete(ctx *Context)
	Patch(ctx *Context)
}

type DefaultService struct {
	Handler http.Handler
}

func (s *DefaultService) defaultHandler(ctx *Context) {
	if s.Handler != nil {
		s.Handler.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	} else {
		ctx.ResponseWriter.WriteHeader(http.StatusNotImplemented)
	}
}

func (s *DefaultService) Get(ctx *Context) {
	s.defaultHandler(ctx)
}

func (s *DefaultService) Post(ctx *Context) {
	s.defaultHandler(ctx)
}

func (s *DefaultService) Put(ctx *Context) {
	s.defaultHandler(ctx)
}

func (s *DefaultService) Delete(ctx *Context) {
	s.defaultHandler(ctx)
}

func (s *DefaultService) Patch(ctx *Context) {
	s.defaultHandler(ctx)
}
