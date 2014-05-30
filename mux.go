package handy

import (
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	Logger *log.Logger
)

type Handy struct {
	mu             sync.RWMutex
	router         *Router
	currentClients int
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
	handy.mu.RLock()
	defer handy.mu.RUnlock()

	handy.currentClients++
	defer func() {
		handy.currentClients--
	}()

	route, err := handy.router.Match(r.URL.Path)
	if err != nil {
		http.Error(rw, "", http.StatusServiceUnavailable)
		return
	}

	constructor := route.Handler
	h := constructor()

	w := &ResponseWriter{ResponseWriter: rw}
	paramsDecoder := ParamCodec{URIParams: route.URIVars}
	paramsDecoder.Decode(w, r, h)

	if codec, ok := h.(Codec); ok {
		codec.Decode(w, r, h)
	}

	if w.Written() {
		return
	}

	executeChain(h.Interceptors(), w, r, h)

	if codec, ok := h.(Codec); ok {
		codec.Encode(w, r, h)
	}
}

func executeChain(is []Interceptor, w *ResponseWriter, r *http.Request, h Handler) {
	if len(is) == 0 {
		call(w, r, h)
		return
	}

	is[0].Before(w, r, h)

	if !w.Written() {
		executeChain(is[1:], w, r, h)
		is[0].After(w, r, h)
	}
}

func call(w http.ResponseWriter, r *http.Request, h Handler) {
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
