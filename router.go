package mux

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrRouteNotFound = errors.New("Path not found")
)

type node struct {
	name             string
	isWildcard       bool
	hasChildWildcard bool
	parent           *node
	children         map[string]*node
	services         []Service
	wildcardName     string
}

type Path struct {
	root    *node
	current *node
}

func NewPath() *Path {
	r := new(Path)
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

func (r *Path) nodeExists(n string) (*node, bool) {
	v, ok := r.current.children[n]
	if !ok && r.current.hasChildWildcard {
		// looking for wildcard
		v, ok = r.current.children[r.current.wildcardName]
	}

	return v, ok
}

func (r *Path) AppendRoute(uri string, s Service) bool {
	appended := false
	for _, v := range strings.Split(uri, "/") {
		if len(v) > 0 {
			if n, ok := r.nodeExists(v); ok {
				r.current = n
				appended = true
				continue
			}

			n := new(node)
			n.children = make(map[string]*node)

			// only one child wildcard per node
			if isWildcard(v) {
				if r.current.hasChildWildcard {
					return false
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
	}

	if r.current != r.root {
		r.current.services = append(r.current.services, s)
		r.current = r.root // reset
	}

	return appended

}

func (n *node) findChild(name string) *node {
	v, ok := n.children[name]
	if !ok && n.hasChildWildcard {
		// looking for wildcard
		v, ok = n.children[n.wildcardName]
	}

	return v
}

type Route struct {
	Route      string
	URIVars    map[string]string
	Services   []Service
	Restricted bool
}

// This method rebuilds a route based on a given URI
func (r *Path) FindRoute(uri string) (Route, error) {
	rt := Route{}
	rt.URIVars = make(map[string]string)

	current := r.current
	for _, v := range strings.Split(uri, "/") {
		v = strings.TrimSpace(v)
		if len(v) > 0 {
			n := current.findChild(v)
			if n == nil {
				return rt, ErrRouteNotFound
			}

			if n.isWildcard {
				rt.URIVars[cleanWildcard(n.name)] = v
			}

			rt.Route = fmt.Sprintf("%s/%s", rt.Route, n.name)
			current = n
		}
	}

	if len(current.services) == 0 {
		return rt, ErrRouteNotFound
	}

	rt.Services = current.services
	return rt, nil
}
