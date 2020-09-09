package handy

import "net/http"

type Handler interface {
	Get() int
	Post() int
	Put() int
	Delete() int
	Patch() int
	Head() int
	SetContext(Context)
}

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	URIVars        URIVars
}

type ProtoHandler struct {
	Context
}

func (h *ProtoHandler) Get() int {
	return http.StatusMethodNotAllowed
}

func (h *ProtoHandler) Post() int {
	return http.StatusMethodNotAllowed
}

func (h *ProtoHandler) Put() int {
	return http.StatusMethodNotAllowed
}

func (h *ProtoHandler) Delete() int {
	return http.StatusMethodNotAllowed
}

func (h *ProtoHandler) Patch() int {
	return http.StatusMethodNotAllowed
}

func (h *ProtoHandler) Head() int {
	return http.StatusMethodNotAllowed
}

func (h *ProtoHandler) SetContext(c Context) {
	h.Context = c
}
