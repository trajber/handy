package handy

import "testing"

func TestAppendRoute(t *testing.T) {
	rt := newRouter()
	h := new(ProtoHandler)
	err := rt.appendRoute("/test/test", func() (Handler, Interceptor) { return h, nil })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.appendRoute("/test", func() (Handler, Interceptor) { return h, nil })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.appendRoute("/test/test", func() (Handler, Interceptor) { return h, nil })
	if err == nil {
		t.Fatal("Appending the same route twice")
	}

	err = rt.appendRoute("/test", func() (Handler, Interceptor) { return h, nil })
	if err == nil {
		t.Fatal("Appending the same route twice")
	}

	err = rt.appendRoute("/test/", func() (Handler, Interceptor) { return h, nil })
	if err == nil {
		t.Fatal("Appending the same route twice", err)
	}
}

func TestAppendWildCard(t *testing.T) {
	rt := newRouter()
	h := new(ProtoHandler)
	err := rt.appendRoute("/test/{x}", func() (Handler, Interceptor) { return h, nil })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.appendRoute("/test/{x}/test", func() (Handler, Interceptor) { return h, nil })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.appendRoute("/test/{x}", func() (Handler, Interceptor) { return h, nil })
	if err == nil {
		t.Fatal("Appending the same route twice")
	}

	err = rt.appendRoute("/test/{x}/test", func() (Handler, Interceptor) { return h, nil })
	if err == nil {
		t.Fatal("Appending the same route twice")
	}
}

func TestAppendInvalidWildCard(t *testing.T) {
	rt := newRouter()
	h := new(ProtoHandler)

	err := rt.appendRoute("/test/{x}", func() (Handler, Interceptor) { return h, nil })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.appendRoute("/test/{y}", func() (Handler, Interceptor) { return h, nil })
	t.Log(err)
	if err == nil {
		t.Fatal("A invalid node was appended", err)
	}
}

func TestFindRoute(t *testing.T) {
	rt := newRouter()
	h := new(ProtoHandler)

	err := rt.appendRoute("/test", func() (Handler, Interceptor) { return h, nil })
	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route := rt.match("/test")
	if route == nil {
		t.Fatal("Cannot find a valid route")
	}

	t.Log(route.URIVars)
}

func TestMatchWithWildcard(t *testing.T) {
	rt := newRouter()
	h := new(ProtoHandler)
	err := rt.appendRoute("/test/{x}", func() (Handler, Interceptor) { return h, nil })

	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route := rt.match("/test/foo")
	if route == nil {
		t.Fatal("Cannot find a valid route")
	}

	t.Log(route.URIVars)
}

func TestAppendSameRoute(t *testing.T) {
	rt := newRouter()
	h := new(ProtoHandler)
	err := rt.appendRoute("/test", func() (Handler, Interceptor) { return h, nil })

	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	err = rt.appendRoute("/test", func() (Handler, Interceptor) { return h, nil })

	if err == nil {
		t.Fatal("Overriting route. This sould not happen.")
	}
}

func TestMultipleWildCards(t *testing.T) {
	rt := newRouter()
	h := new(ProtoHandler)
	err := rt.appendRoute("/test/{x}/{y}", func() (Handler, Interceptor) { return h, nil })

	if err != nil {
		t.Fatal("Cannot append a valid route", err)
	}

	route := rt.match("/test/foo/bar")
	if route == nil {
		t.Fatal("Cannot find a valid route;", err)
	}

	t.Log(route.URIVars)
}
