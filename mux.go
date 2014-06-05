package handy

import (
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
)

var (
	Logger *log.Logger
)

type Handy struct {
	mu             sync.RWMutex
	router         *Router
	currentClients int32
	CountClients   bool
}

type HandyFunc func() Handler

func NewHandy() *Handy {
	handy := new(Handy)
	handy.router = NewRouter()
	Logger = log.New(os.Stdout, "[handy] ", 0)
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

	route, err := handy.router.Match(r.URL.Path)
	if err != nil {
		http.Error(rw, "", http.StatusServiceUnavailable)
		return
	}

	h := route.Handler()

	w := &ResponseWriter{ResponseWriter: rw}
	paramsDecoder := newParamDecoder(h, route.URIVars)
	paramsDecoder.Decode(w, r)

	interceptors := h.Interceptors()
	for k, interceptor := range interceptors {
		interceptor.Before(w, r)
		if !w.Written() {
			continue
		}

		// if something was written... pop-out all executed interceptors
		// and execute them in reverse order calling After method.
		for rev := k; rev >= 0; rev-- {
			interceptors[rev].After(w, r)
		}

		return
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

	// executing all After interceptors in reverse order
	for k, _ := range interceptors {
		interceptors[len(interceptors)-1-k].After(w, r)
	}
}
