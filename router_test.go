package handy

import "testing"

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

func TestAppendInvalidWildCard(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)

	err := rt.AppendRoute("/test/{x}", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.AppendRoute("/test/{y}", func() Handler { return h })
	t.Log(err)
	if err == nil {
		t.Fatal("A invalid node was appended", err)
	}
}

func TestDontFindRoute(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)

	err := rt.AppendRoute("/test", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route, err := rt.Match("/test2")
	if err == nil {
		t.Fatal("Got a match for an invalid route")
	}

	route, err = rt.Match("/test/test2")
	if err == nil {
		t.Fatal("Got a match for an invalid route")
	}

	t.Log(route.URIVars)
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

	t.Log(route.URIVars)
}

func TestFindRouteWithEmptyParameter(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)

	err := rt.AppendRoute("/test/{param}", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route, err := rt.Match("/test/")
	if err != nil {
		t.Fatal("Cannot find a valid route;", err)
	}

	t.Log(route.URIVars)
}

func TestFindRouteFallback(t *testing.T) {
	rt := NewRouter()
	h := new(DefaultHandler)

	err := rt.AppendRoute("/test/{param}", func() Handler { return h })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route, err := rt.Match("/test/value/another-test")
	if err != nil {
		t.Fatal("Cannot find a valid route;", err)
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
