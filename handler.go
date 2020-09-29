package handy

import "net/http"

// Handler is the interface that any handler implements.
//
// When an HTTP request arrives, the framework calls the corresponding
// handler's method based on the HTTP request. The integer returned by them are
// then injected in the After method of each registered interceptor. Each
// interceptor may choose what to do with this value. Although it's expected
// the return value be the HTTP status code of the response, the framework
// imposes no meaning to it and in particular does not write it to the
// response's header. It's up to each interceptor decide what to do with it.
// See the documentation of Interceptor for more details.
//
// All handlers must also embed ProtoHandler to make them compatible with the
// framework.
type Handler interface {
	Get() int
	Post() int
	Put() int
	Delete() int
	Patch() int
	Head() int

	// SetContext is used internally by the framework to set Context
	// information on each handler. It's not meant to be called by the user,
	// but it's exported in order to be used to inject mock data during your
	// tests.
	SetContext(Context)
}

// Context is automatically embedded in each interceptor and handler, giving
// access to request information and to the ResponseWriter object. For a
// handler, it's rarely used, since it's expected the request and response data
// to be managed by interceptors in a type safe manner. But it's available for
// when you face specific situations for which the work done by interceptors is
// not enough.
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	// URIVars gives you access to the variables declared in the handler's route.
	URIVars map[string]string
}

// ProtoHandler is a prototype implementation for a handler. It must be embedded
// in all of your handlers to make them compatible with Handy.
type ProtoHandler struct {
	Context
}

// Get is a default implementation that must be overwritten for
// handling GET requests. This default implementation always answers
// http.StatusMethodNotAllowed
func (h *ProtoHandler) Get() int {
	return http.StatusMethodNotAllowed
}

// Post is a default implementation that must be overwritten for
// handling POST requests. This default implementation always answers
// http.StatusMethodNotAllowed
func (h *ProtoHandler) Post() int {
	return http.StatusMethodNotAllowed
}

// Put is a default implementation that must be overwritten for
// handling PUT requests. This default implementation always answers
// http.StatusMethodNotAllowed
func (h *ProtoHandler) Put() int {
	return http.StatusMethodNotAllowed
}

// Delete is a default implementation that must be overwritten for
// handling DELETE requests. This default implementation always answers
// http.StatusMethodNotAllowed
func (h *ProtoHandler) Delete() int {
	return http.StatusMethodNotAllowed
}

// Patch is a default implementation that must be overwritten for
// handling PATCH requests. This default implementation always answers
// http.StatusMethodNotAllowed
func (h *ProtoHandler) Patch() int {
	return http.StatusMethodNotAllowed
}

// Head is a default implementation that must be overwritten for
// handling HEAD requests. This default implementation always answers
// http.StatusMethodNotAllowed
func (h *ProtoHandler) Head() int {
	return http.StatusMethodNotAllowed
}

func (h *ProtoHandler) SetContext(c Context) {
	h.Context = c
}
