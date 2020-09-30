// Handy is a fast and simple HTTP multiplexer. It favors composition by
// introducing the concept of interceptors (an implementation of the decorator
// pattern), allowing you to reduce the logic of your handler to a minimum
// specific part, whereas the common logic shared by several handlers is
// implemented in many composable interceptors.
//
// In its most basic form, it allows the logic of your handler to be split by HTTP method:
//
//     func main() {
//         server := handy.New()
//         server.Handle("/user/{username}", func() (handy.Handler, handy.Interceptor) {
//             return &userHandler{}, nil
//         })
//
//         http.ListenAndServe(":8181", server)
//     }
//
//     type userHandler struct {
//         handy.BaseHandler
//     }
//
//     func (h *userHandler) Get() int {
//         username := h.URIVars["username"]
//         response := ...
//         h.ResponseWriter.Write(...)
//         return http.StatusOK
//     }
//
// The true power comes when you plug interceptors into the pipeline:
//
//     func main() {
//         server := handy.New()
//         server.Handle("population/{city}/{year}", func() (handy.Handler, handy.Interceptor) {
//             handler := new(userHandler)
//             introspector := interceptor.NewIntrospector(nil, handler)
//             uriVars := interceptor.NewURIVars(introspector)
//             codec := interceptor.NewJSONCodec(uriVars)
//             return handler, codec
//         })
//
//         http.ListenAndServe(":8181", server)
//     }
//
//     type userHandler struct {
//         handy.BaseHandler
//
//         City string `urivar:"city"`
//         Year int `urivar:"year"`
//         Response Statistics `response:"get"`
//     }
//
//     func (h *userHandler) Get() int {
//         statistics := populationByCityAndYear(h.City, h.Year)
//         h.Response = statistics
//         return http.StatusOK
//     }
//
// As you can see, interceptors can automatically parse URI variables and handle
// marshaling and unmarshaling of requests and responses. And you can also
// write your custom interceptors to do all sort of stuff, like database
// transaction management and data validation.
package handy

import (
	"fmt"
	"net/http"
	"sync"
)

// CatchAllHandler, if defined, is called whenever a request is made for a route
// without any registered handler. You can use it, for instance, to send custom
// 404 responses.
var CatchAllHandler http.Handler

// Handy is the multiplexer of the framework.
type Handy struct {
	// Recover is called whenever a handler panics handling a request.
	//
	// If a handler panics at some point, the framework recovers from the
	// panicked goroutine, answers with a http.StatusInternalServerError and,
	// if the Recover function is set, it's called with any error value
	// retrieved from the call of panic.
	Recover func(err interface{})

	mu     sync.RWMutex
	router *router
}

// New returns a new Handy multiplexer.
func New() *Handy {
	handy := new(Handy)
	handy.router = newRouter()
	return handy
}

// Handle registers a handler to be called whenever a route matches.
//
// The route is in the format /a/{var1}/b/{var2}, where var1 and var2 are
// variables that match anything and whose values are available to be inspected
// in the URIVars field present in each handler and interceptor.
// The handler argument is a constructor returning both a Handler and the
// Interceptor that decorates it.
func (handy *Handy) Handle(route string, handler func() (Handler, Interceptor)) {
	handy.mu.Lock()
	defer handy.mu.Unlock()

	if err := handy.router.appendRoute(route, handler); err != nil {
		panic(fmt.Sprintf("cannot append route “%s”: %v", route, err.Error()))
	}
}

// ServeHTTP makes Handy adheres to the standard http.Handler interface. As such, it can be used whenever a http.Handler is expected.
func (handy *Handy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handy.mu.RLock()
	defer handy.mu.RUnlock()

	defer func() {
		if r := recover(); r != nil {
			if handy.Recover != nil {
				handy.Recover(r)
			}
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}()

	route := handy.router.match(request.URL.Path)

	if route == nil {
		if CatchAllHandler != nil {
			CatchAllHandler.ServeHTTP(writer, request)
		} else {
			writer.WriteHeader(http.StatusNotFound)
		}
		return
	}

	c := Context{ResponseWriter: writer, Request: request, URIVars: route.URIVars}
	handler, interceptors := route.Handler()
	handler.SetContext(c)
	interceptors.SetContext(c)
	chain := buildChain(interceptors)

	var status int

	// chain is in reverse order
	for k := len(chain) - 1; k >= 0; k-- {
		interceptor := chain[k]
		status = interceptor.Before()

		// If the interceptor reports some status, interrupt the chain
		if status != 0 {
			chain = chain[k:]
			goto write
		}
	}

	switch request.Method {
	case http.MethodGet:
		status = handler.Get()
	case http.MethodPost:
		status = handler.Post()
	case http.MethodPut:
		status = handler.Put()
	case http.MethodDelete:
		status = handler.Delete()
	case http.MethodPatch:
		status = handler.Patch()
	case http.MethodHead:
		status = handler.Head()
	default:
		status = http.StatusMethodNotAllowed
	}

write:
	// executing all interceptors' After methods, in reverse order
	for _, interceptor := range chain {
		s := interceptor.After(status)

		if s != 0 {
			status = s
		}
	}
}
