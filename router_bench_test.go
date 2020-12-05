package handy

import "testing"

func BenchmarkFindRoute(b *testing.B) {
	rt := newRouter()
	h := new(BaseHandler)
	err := rt.appendRoute("/test/{x}", func() (Handler, Interceptor) {
		return h, nil
	})

	if err != nil {
		b.Fatal("Cannot append a valid route", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		route := rt.match("/test/foo")
		if route == nil {
			b.Fatal("Cannot find a valid route")
		}
	}
}
