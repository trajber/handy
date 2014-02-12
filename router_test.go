package handy

import (
	"testing"
)

func TestAppendRoute(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	err := rt.AppendRoute("/test", h)
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}
}

func TestAppendWildCard(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	err := rt.AppendRoute("/test/{x}", h)
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}
}

func TestFindRoute(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	t.Log(h)
	err := rt.AppendRoute("/test", h)

	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route, err := rt.FindRoute("/test")
	if err != nil {
		t.Fatal("Cannot find a valid route;", err)
	}

	t.Log(route.URIVars)
}

func TestFindRouteWithWildcard(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	t.Log(h)
	err := rt.AppendRoute("/test/{x}", h)

	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route, err := rt.FindRoute("/test/foo")
	if err != nil {
		t.Fatal("Cannot find a valid route;", err)
	}

	t.Log(route.URIVars)
}

func TestAppendSameRoute(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	t.Log(h)
	err := rt.AppendRoute("/test", h)

	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.AppendRoute("/test", h)

	if err == nil {
		t.Fatal("Overriting route. This sould not happen.")
	}
}

func TestMultipleWildCards(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	t.Log(h)
	err := rt.AppendRoute("/test/{x}/{y}", h)

	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route, err := rt.FindRoute("/test/foo/bar")
	if err != nil {
		t.Fatal("Cannot find a valid route;", err)
	}

	t.Log(route.URIVars)
}
