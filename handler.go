package handy

import "net/http"

type Handler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
	Put(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
	Patch(http.ResponseWriter, *http.Request)
	Decode(http.ResponseWriter, *http.Request, Handler)
	Encode(http.ResponseWriter, *http.Request, Handler)
	Interceptors() InterceptorChain
}

type DefaultHandler struct {
	http.Handler
	NopCodec
	NopInterceptorChain
}

func (s *DefaultHandler) defaultHandler(w http.ResponseWriter, r *http.Request) {
	if s.Handler != nil {
		s.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (s *DefaultHandler) Get(w http.ResponseWriter, r *http.Request) {
	s.defaultHandler(w, r)
}

func (s *DefaultHandler) Post(w http.ResponseWriter, r *http.Request) {
	s.defaultHandler(w, r)
}

func (s *DefaultHandler) Put(w http.ResponseWriter, r *http.Request) {
	s.defaultHandler(w, r)
}

func (s *DefaultHandler) Delete(w http.ResponseWriter, r *http.Request) {
	s.defaultHandler(w, r)
}

func (s *DefaultHandler) Patch(w http.ResponseWriter, r *http.Request) {
	s.defaultHandler(w, r)
}
