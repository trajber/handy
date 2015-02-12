package handy

import (
	"net/http"
	"sync"
	"sync/atomic"
)

var (
	ErrorFunc func(e error)
)

type Handy struct {
	mu             sync.RWMutex
	router         *Router
	currentClients int32
	CountClients   bool
	Recover        func(interface{})
}

type HandyFunc func() Handler

func NewHandy() *Handy {
	handy := new(Handy)
	handy.router = NewRouter()
	return handy
}

func (handy *Handy) Handle(pattern string, h HandyFunc) {
	handy.mu.Lock()
	defer handy.mu.Unlock()

	if err := handy.router.AppendRoute(pattern, h); err != nil {
		panic("Cannot append route;" + err.Error())
	}
}

func (handy *Handy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}()

	route, err := handy.router.Match(r.URL.Path)
	if err != nil {
		// http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.5
		// The server has not found anything matching the Request-URI. No indication is given of whether
		// the condition is temporary or permanent.
		http.NotFound(rw, r)
		return
	}

	h := route.Handler()

	w := NewBufferedResponseWriter(rw)

	paramsDecoder := newParamDecoder(h, route.URIVars)
	paramsDecoder.Decode(w, r)

	interceptors := h.Interceptors()
	for k, interceptor := range interceptors {
		interceptor.Before(w, r)
		// if something was written we need to stop the execution
		if w.status > 0 || (w.flushed || w.Body.Len() > 0) {
			interceptors = interceptors[:k+1]
			goto write
		}
	}

	switch r.Method {
	case "GET":
		h.Get(w, r)
	case "POST":
		h.Post(w, r)
	case "PUT":
		h.Put(w, r)
	case "DELETE":
		h.Delete(w, r)
	case "PATCH":
		h.Patch(w, r)
	case "HEAD":
		h.Head(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}

write:
	// executing all After interceptors in reverse order
	for k, _ := range interceptors {
		interceptors[len(interceptors)-1-k].After(w, r)
	}

	w.Flush()
}
