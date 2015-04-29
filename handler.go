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
	setRequestInfo(w http.ResponseWriter, r *http.Request, u URIVars)
}

type DefaultHandler struct {
	http.Handler
	NopInterceptorChain

	response http.ResponseWriter
	request  *http.Request
	uriVars  URIVars
}

func (d *DefaultHandler) handle() int {
	if d.Handler != nil {
		d.ServeHTTP(d.response, d.request)
		return http.StatusOK
	} else {
		d.response.WriteHeader(http.StatusMethodNotAllowed)
		return http.StatusMethodNotAllowed
	}
}

func (d *DefaultHandler) Get() int {
	return d.handle()
}

func (d *DefaultHandler) Post() int {
	return d.handle()
}

func (d *DefaultHandler) Put() int {
	return d.handle()
}

func (d *DefaultHandler) Delete() int {
	return d.handle()
}

func (d *DefaultHandler) Patch() int {
	return d.handle()
}

func (d *DefaultHandler) Head() int {
	return d.handle()
}

func (d *DefaultHandler) ResponseWriter() http.ResponseWriter {
	return d.response
}

func (d *DefaultHandler) Req() *http.Request {
	return d.request
}

func (d *DefaultHandler) URIVars() URIVars {
	return d.uriVars
}

func (d *DefaultHandler) setRequestInfo(w http.ResponseWriter, r *http.Request, u URIVars) {
	*d = DefaultHandler{response: w, request: r, uriVars: u}
}
