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
	setRequestInfo(w http.ResponseWriter, r *http.Request, uri string, u URIVars)
}

type DefaultHandler struct {
	NopInterceptorChain

	response http.ResponseWriter
	request  *http.Request
	uri      string
	uriVars  URIVars
}

func (d *DefaultHandler) Get() int {
	return http.StatusMethodNotAllowed
}

func (d *DefaultHandler) Post() int {
	return http.StatusMethodNotAllowed
}

func (d *DefaultHandler) Put() int {
	return http.StatusMethodNotAllowed
}

func (d *DefaultHandler) Delete() int {
	return http.StatusMethodNotAllowed
}

func (d *DefaultHandler) Patch() int {
	return http.StatusMethodNotAllowed
}

func (d *DefaultHandler) Head() int {
	return http.StatusMethodNotAllowed
}

func (d *DefaultHandler) ResponseWriter() http.ResponseWriter {
	return d.response
}

func (d *DefaultHandler) Req() *http.Request {
	return d.request
}

func (d *DefaultHandler) URI() string {
	return d.uri
}

func (d *DefaultHandler) URIVars() URIVars {
	return d.uriVars
}

func (d *DefaultHandler) setRequestInfo(w http.ResponseWriter, r *http.Request, uri string, u URIVars) {
	*d = DefaultHandler{response: w, request: r, uri: uri, uriVars: u}
}
