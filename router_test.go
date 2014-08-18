package handy

import (
	"testing"
)

func TestAppendRoute(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	err := rt.AppendRoute("/test/test", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.AppendRoute("/test", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.AppendRoute("/test/test", func() Handler { return h })
	if err == nil {
		t.Fatal("Appending the same route twice")
	}

	err = rt.AppendRoute("/test", func() Handler { return h })
	if err == nil {
		t.Fatal("Appending the same route twice")
	}

	err = rt.AppendRoute("/test/", func() Handler { return h })
	if err == nil {
		t.Fatal("Appending the same route twice", err)
	}

	err = rt.AppendRoute("/", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append root", err)
	}
}

func TestAppendWildCard(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	err := rt.AppendRoute("/test/{x}", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.AppendRoute("/test/{x}/test", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.AppendRoute("/test/{x}", func() Handler { return h })
	if err == nil {
		t.Fatal("Appending the same route twice")
	}

	err = rt.AppendRoute("/test/{x}/test", func() Handler { return h })
	if err == nil {
		t.Fatal("Appending the same route twice")
	}
}

func TestFindRoute(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)

	err := rt.AppendRoute("/test", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route, err := rt.Match("/test")
	if err != nil {
		t.Fatal("Cannot find a valid route;", err)
	}

	route, err = rt.Match("/test/test")
	if err != nil {
		t.Fatal("Not falling back to parent handler;", err)
	}

	err = rt.AppendRoute("/another/{param}", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route, err = rt.Match("/another/{param}/test")
	if err != nil {
		t.Fatal("Not falling back to parent handler;", err)
	}

	t.Log(route.URIVars)
}

func TestMatchWithWildcard(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	err := rt.AppendRoute("/test/{x}", func() Handler { return h })

	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route, err := rt.Match("/test/foo")
	if err != nil {
		t.Fatal("Cannot find a valid route;", err)
	}

	t.Log(route.URIVars)
}

func TestAppendSameRoute(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	err := rt.AppendRoute("/test", func() Handler { return h })

	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.AppendRoute("/test", func() Handler { return h })

	if err == nil {
		t.Fatal("Overriting route. This sould not happen.")
	}
}

func TestMultipleWildCards(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)
	err := rt.AppendRoute("/test/{x}/{y}", func() Handler { return h })

	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route, err := rt.Match("/test/foo/bar")
	if err != nil {
		t.Fatal("Cannot find a valid route;", err)
	}

	t.Log(route.URIVars)
}
