package handy

import "net/http"

type Handler interface {
	Get(ctx *Context) (int, error)
	Post(ctx *Context) (int, error)
	Put(ctx *Context) (int, error)
	Delete(ctx *Context) (int, error)
	Patch(ctx *Context) (int, error)
	Decoder(f Interceptor)
	getDecoder() Interceptor
	Encoder(f Interceptor)
	getEncoder() Interceptor
	After(f ...Interceptor)
	getChain() []Interceptor
}

type DefaultHandler struct {
	http.Handler
	encoder         Interceptor
	decoder         Interceptor
	getInterceptors []Interceptor
}

func (s *DefaultHandler) defaultHandler(ctx *Context) (int, error) {
	if s.Handler != nil {
		s.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	} else {
		ctx.ResponseWriter.WriteHeader(http.StatusNotImplemented)
	}

	return http.StatusNotImplemented, nil
}

func (s *DefaultHandler) After(interceptors ...Interceptor) {
	s.getInterceptors = make([]Interceptor, 0)
	s.getInterceptors = append(s.getInterceptors, interceptors...)
}

func (s *DefaultHandler) Decoder(dec Interceptor) {
	s.decoder = dec
}

func (s *DefaultHandler) getDecoder() Interceptor {
	return s.decoder
}

func (s *DefaultHandler) Encoder(enc Interceptor) {
	s.encoder = enc
}

func (s *DefaultHandler) getEncoder() Interceptor {
	return s.encoder
}

func (s *DefaultHandler) getChain() []Interceptor {
	return s.getInterceptors
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
