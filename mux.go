package handy

import (
	"net/http"
	"sync"
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
	return handy
}

func (handy *Handy) HandleService(pattern string, h HandyFunc) {
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

	h.Decode(w, r, h)

	if w.Written() {
		return
	}

	for _, i := range h.Before() {
		i.Intercept(w, r, h)
		if w.Written() {
			return
		}
	}

	if w.Written() {
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
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}

	if w.Written() {
		return
	}

	for _, i := range h.After() {
		i.Intercept(w, r, h)
		if w.Written() {
			return
		}
	}

	h.Encode(w, r, h)
}

func (handy *Handy) Handle(pattern string, handler http.Handler) {
	panic("unsupported")
	// here we have to create a Handler with all verbs calling 'handler'
	// handy.mu.Lock()
	// defer handy.mu.Unlock()
}

func (handy *Handy) HandleFunc(pattern string,
	handler func(http.ResponseWriter, *http.Request)) {
	handy.Handle(pattern, http.HandlerFunc(handler))
}
