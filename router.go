package handy

import (
	"errors"
	"strings"
)

var (
	ErrRouteNotFound      = errors.New("Router not found")
	ErrRouteAlreadyExists = errors.New("Route already exists")
	ErrCannotAppendRoute  = errors.New("Cannot append route")
	ErrOnlyOneWildcard    = errors.New("Only one wildcard is allowed in this level")
)

type node struct {
	name             string
	handler          Constructor
	isWildcard       bool
	hasChildWildcard bool
	parent           *node
	children         map[string]*node
	wildcardName     string
}

type Router struct {
	root    *node
	current *node
}

func NewRouter() *Router {
	r := new(Router)
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

func (r *Router) nodeExists(n string) (*node, bool) {
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

func (r *Router) AppendRoute(uri string, h Constructor) error {
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
			return ErrCannotAppendRoute
		}

		if n, ok := r.nodeExists(v); ok {
			if i == len(tokens)-1 && n.handler != nil {
				return ErrRouteAlreadyExists

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
				return ErrOnlyOneWildcard
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
		return ErrCannotAppendRoute
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

type URIVars map[string]string

type RouteMatch struct {
	URIVars URIVars
	Handler Constructor
}

// This method rebuilds a route based on a given URI
func (r *Router) Match(uri string) (*RouteMatch, error) {
	rt := new(RouteMatch)
	rt.URIVars = make(URIVars)

	current := r.current
	uri = strings.TrimSpace(uri)
	parts := strings.Split(uri, "/")

	for i, v := range parts {
		// ignore first empty value (before initial slash)
		if i == 0 && v == "" {
			continue
		}

		n := current.findChild(v)
		if n == nil {
			if !current.isWildcard {
				return rt, ErrRouteNotFound
			}

			// when we cannot find the specific route, fallback to the last handler
			// that we found. The URI parts that did not match will be concataned to
			// the URI variable when the last handler is a wildcard
			rt.URIVars[cleanWildcard(current.name)] += "/" + strings.Join(parts[i:], "/")
			break
		}

		if n.isWildcard {
			rt.URIVars[cleanWildcard(n.name)] = v
		}

		current = n
	}

	if current.handler == nil {
		return rt, ErrRouteNotFound
	}

	rt.Handler = current.handler
	return rt, nil
}
