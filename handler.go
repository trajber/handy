package handy

import "net/http"

type Handler interface {
	Get() int
	Post() int
	Put() int
	Delete() int
	Patch() int
	Head() int
	Interceptors() InterceptorChain
}

func BuildDefaultHandler(w http.ResponseWriter, r *http.Request, u URIVars) DefaultHandler {
	return DefaultHandler{response: w, request: r, uriVars: u}
}

type DefaultHandler struct {
	http.Handler
	NopInterceptorChain

	response http.ResponseWriter
	request  *http.Request
	uriVars  URIVars
}

func (s *DefaultHandler) handle() int {
	if s.Handler != nil {
		s.ServeHTTP(s.response, s.request)
		return http.StatusOK
	} else {
		s.response.WriteHeader(http.StatusMethodNotAllowed)
		return http.StatusMethodNotAllowed
	}
}

func (s *DefaultHandler) Get() int {
	return s.handle()
}

func (s *DefaultHandler) Post() int {
	return s.handle()
}

func (s *DefaultHandler) Put() int {
	return s.handle()
}

func (s *DefaultHandler) Delete() int {
	return s.handle()
}

func (s *DefaultHandler) Patch() int {
	return s.handle()
}

func (s *DefaultHandler) Head() int {
	return s.handle()
}

func (s *DefaultHandler) ResponseWriter() http.ResponseWriter {
	return s.response
}

func (s *DefaultHandler) Request() *http.Request {
	return s.request
}

func (s *DefaultHandler) URIVars() URIVars {
	return s.uriVars
}
