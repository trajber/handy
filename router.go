package handy

import (
	"errors"
	"strings"
)

var (
	errRouteAlreadyExists = errors.New("route already exists")
	errCannotAppendRoute  = errors.New("cannot append route")
	errOnlyOneWildcard    = errors.New("only one wildcard is allowed in this level")
)

type node struct {
	name             string
	handler          func() (Handler, Interceptor)
	isWildcard       bool
	hasChildWildcard bool
	parent           *node
	children         map[string]*node
	wildcardName     string
}

type router struct {
	root    *node
	current *node
}

func newRouter() *router {
	r := new(router)
	root := new(node)
	root.children = make(map[string]*node)
	r.root = root
	r.current = r.root
	return r
}

func isWildcard(l string) bool {
	return l[0] == '{' && l[len(l)-1] == '}'
}

func cleanWildcard(l string) string {
	return l[1 : len(l)-1]
}

func (r *router) nodeExists(n string) (*node, bool) {
	v, ok := r.current.children[n]
	if !ok && r.current.hasChildWildcard {
		if isWildcard(n) {
			n = cleanWildcard(n)
			// looking for wildcard with the same name
			v, ok = r.current.children[n]
		}
	}

	return v, ok
}

func (r *router) appendRoute(uri string, h func() (Handler, Interceptor)) error {
	uri = strings.TrimSpace(uri)

	// Make sure we are not appending the root ("/"), otherwise remove final slash
	if len(uri) > 1 && uri[len(uri)-1] == '/' {
		uri = uri[:len(uri)-1]
	}

	// Should end at root node
	defer func() {
		r.current = r.root
	}()

	appended := false
	tokens := strings.Split(uri, "/")
	for i, v := range tokens {
		if v == "" {
			continue
		}

		if r.current.hasChildWildcard && !isWildcard(v) {
			return errCannotAppendRoute
		}

		if n, ok := r.nodeExists(v); ok {
			if i == len(tokens)-1 && n.handler != nil {
				return errRouteAlreadyExists

			} else if i == len(tokens)-1 {
				n.handler = h
				return nil
			}

			r.current = n
			appended = true
			continue
		}

		n := new(node)
		n.children = make(map[string]*node)

		// only one child wildcard per node
		if isWildcard(v) {
			if r.current.hasChildWildcard {
				return errOnlyOneWildcard
			}

			n.isWildcard = true
			r.current.wildcardName = v
			r.current.hasChildWildcard = true
		}

		n.name = v
		n.parent = r.current
		r.current.children[n.name] = n
		r.current = n
		appended = true
	}

	if r.current != r.root {
		r.current.handler = h
	}

	if appended == false {
		return errCannotAppendRoute
	}

	return nil

}

func (n *node) findChild(name string) *node {
	v, ok := n.children[name]
	if !ok && n.hasChildWildcard {
		// looking for wildcard
		v = n.children[n.wildcardName]
	}

	return v
}

type routeMatch struct {
	URIVars map[string]string
	Handler func() (Handler, Interceptor)
}

// This method rebuilds a route based on a given URI
func (r *router) match(uri string) *routeMatch {
	rt := new(routeMatch)
	rt.URIVars = make(map[string]string)

	current := r.current
	uri = strings.TrimSpace(uri)
	for _, v := range strings.Split(uri, "/") {
		if v == "" {
			continue
		}

		n := current.findChild(v)
		if n == nil {
			return nil
		}

		if n.isWildcard {
			rt.URIVars[cleanWildcard(n.name)] = v
		}

		current = n
	}

	if current.handler == nil {
		return nil
	}

	rt.Handler = current.handler
	return rt
}
