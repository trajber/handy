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
	countClients   bool
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

func (handy *Handy) CountClients() {
	handy.countClients = true
}

func (handy *Handy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if handy.countClients {
		atomic.AddInt32(&handy.currentClients, 1)
		defer atomic.AddInt32(&handy.currentClients, -1)
	}

	route, err := handy.router.Match(r.URL.Path)
	if err != nil {
		http.Error(rw, "", http.StatusServiceUnavailable)
		return
	}

	h := route.Handler()

	w := &ResponseWriter{ResponseWriter: rw}
	paramsDecoder := newParamDecoder(h, route.URIVars)
	paramsDecoder.Decode(w, r)

	executeChain(h.Interceptors(), h, w, r)
}

func executeChain(is []Interceptor, h Handler, w *ResponseWriter, r *http.Request) {
	if len(is) == 0 {
		call(h, w, r)
		return
	}

	is[0].Before(w, r)

	if !w.Written() {
		executeChain(is[1:], h, w, r)
		is[0].After(w, r)
	}
}

func call(h Handler, w http.ResponseWriter, r *http.Request) {
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
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
