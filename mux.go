package handy

import (
	"net/http"
	"sync"
)

type Handy struct {
	mu                sync.RWMutex
	router            *Router
	beforeFilter      Filter
	afterFilter       Filter
	handleFilterError FilterError
	currentClients    int
}

type Filter func(ctx *Context) error
type FilterError func(ctx *Context, err error)
type HandyFunc func(ctx *Context) Handler

func NewHandy() *Handy {
	handy := new(Handy)
	handy.router = NewRouter()
	return handy
}

func (handy *Handy) Handle(pattern string, handler http.Handler) {
	// handy.mu.Lock()
	// defer handy.mu.Unlock()

	// s := new(DefaultHandler)
	// s.Handler = handler
	// if err := handy.router.AppendRoute(pattern, s); err != nil {
	// 	panic("Cannot append route:" + err.Error())
	// }
}

func (handy *Handy) HandleFunc(pattern string,
	handler func(http.ResponseWriter, *http.Request)) {
	handy.Handle(pattern, http.HandlerFunc(handler))
}

func (handy *Handy) HandleService(pattern string, h HandyFunc) {
	handy.mu.Lock()
	defer handy.mu.Unlock()

	if err := handy.router.AppendRoute(pattern, h); err != nil {
		panic("Cannot append route;" + err.Error())
	}
}

func (handy *Handy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handy.mu.RLock()
	defer handy.mu.RUnlock()

	handy.currentClients++
	defer func() {
		handy.currentClients--
	}()

	route, err := handy.router.Match(r.URL.Path)
	if err != nil {
		http.Error(w, "", http.StatusServiceUnavailable)
		return
	}

	ctx := newContext()
	ctx.Request = r
	ctx.ResponseWriter = w
	ctx.vars = route.URIVars

	if handy.beforeFilter != nil {
		if err := handy.beforeFilter(ctx); err != nil {
			if handy.handleFilterError != nil {
				handy.handleFilterError(ctx, err)
			}
			return
		}
	}

	f := route.Handler
	h := f(ctx)

	h.Decode(ctx, h)

	for _, i := range h.Interceptors() {
		i.Intercept(ctx, h)
	}

	switch r.Method {
	case "GET":
		h.Get(ctx)
	case "POST":
		h.Post(ctx)
	case "PUT":
		h.Put(ctx)
	case "DELETE":
		h.Delete(ctx)
	case "PATCH":
		h.Patch(ctx)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}

	h.Encode(ctx, h)

	if handy.afterFilter != nil {
		if err := handy.afterFilter(ctx); err != nil {
			if handy.handleFilterError != nil {
				handy.handleFilterError(ctx, err)
			}
		}
	}
}

func (handy *Handy) BeforeFilter(f Filter) {
	handy.beforeFilter = f
}

func (handy *Handy) AfterFilter(f Filter) {
	handy.afterFilter = f
}

func (handy *Handy) HandleFilterError(f FilterError) {
	handy.handleFilterError = f
}
