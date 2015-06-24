package handy

import (
	"net/http"
	"sync"
	"sync/atomic"
)

var (
	ErrorFunc = func(e error) {}
)

type Handy struct {
	mu             sync.RWMutex
	router         *Router
	currentClients int32
	CountClients   bool
	Recover        func(interface{})
}

type Constructor func() Handler

func SetHandlerInfo(h Handler, w http.ResponseWriter, r *http.Request, u URIVars) {
	h.setRequestInfo(w, r, u)
}

func NewHandy() *Handy {
	handy := new(Handy)
	handy.router = NewRouter()
	return handy
}

func (handy *Handy) Handle(pattern string, h Constructor) {
	handy.mu.Lock()
	defer handy.mu.Unlock()

	if err := handy.router.AppendRoute(pattern, h); err != nil {
		panic("Cannot append route;" + err.Error())
	}
}

func (handy *Handy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	route, err := handy.router.Match(r.URL.Path)

	if err != nil {
		// http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.5
		// The server has not found anything matching the Request-URI. No
		// indication is given of whether the condition is temporary or
		// permanent.
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h := route.Handler()
	SetHandlerInfo(h, w, r, route.URIVars)
	interceptors := h.Interceptors()
	var status int

	for k, interceptor := range interceptors {
		status = interceptor.Before()
		// If the interceptor reported some status, interrupt the chain
		if status != 0 {
			interceptors = interceptors[:k+1]
			goto write
		}
	}

	switch r.Method {
	case "GET":
		status = h.Get()
	case "POST":
		status = h.Post()
	case "PUT":
		status = h.Put()
	case "DELETE":
		status = h.Delete()
	case "PATCH":
		status = h.Patch()
	case "HEAD":
		status = h.Head()
	default:
		status = http.StatusMethodNotAllowed
	}

write:
	// executing all After interceptors in reverse order
	for k := len(interceptors) - 1; k >= 0; k-- {
		s := interceptors[k].After(status)

		if s != 0 {
			status = s
		}
	}
}
