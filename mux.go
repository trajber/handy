// Handy is a fast and simple HTTP multiplexer. It favors composition by
// introducing the concept of interceptors, allowing you to reduce the logic
// of your handler to a minimum specific part, whereas the common logic shared
// by several handlers is implemented in many composable interceptors.
package handy

import (
	"net/http"
	"sync"
	"sync/atomic"
)

var CatchAllHandler http.Handler

type Handy struct {
	mu             sync.RWMutex
	router         *router
	currentClients int32
	CountClients   bool
	Recover        func(interface{})
}

func New() *Handy {
	handy := new(Handy)
	handy.router = newRouter()
	return handy
}

func (handy *Handy) Handle(pattern string, h HandlerConstructor) {
	handy.mu.Lock()
	defer handy.mu.Unlock()

	if err := handy.router.appendRoute(pattern, h); err != nil {
		panic("Cannot append route;" + err.Error())
	}
}

func (handy *Handy) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if handy.CountClients {
		atomic.AddInt32(&handy.currentClients, 1)
		defer atomic.AddInt32(&handy.currentClients, -1)
	}

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

	c := Context{
		ResponseWriter: writer,
		Request:        request,
		URIVars:        route.URIVars,
	}
	handler, interceptors := route.Handler()
	handler.SetContext(c)
	interceptors.SetContext(c)
	chain := buildChain(interceptors)

	var status int

	// chain is in reverse order
	for k := len(chain) - 1; k >= 0; k-- {
		interceptor := chain[k]
		status = interceptor.Before()

		// If the interceptor reported some status, interrupt the chain
		if status != 0 {
			chain = chain[k:]
			goto write
		}
	}

	switch request.Method {
	case "GET":
		status = handler.Get()
	case "POST":
		status = handler.Post()
	case "PUT":
		status = handler.Put()
	case "DELETE":
		status = handler.Delete()
	case "PATCH":
		status = handler.Patch()
	case "HEAD":
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
