package handy

import "testing"

func BenchmarkFindRoute(b *testing.B) {
	rt := NewRouter()
	h := new(ProtoHandler)
	err := rt.AppendRoute("/test/{x}", func() (Handler, Interceptor) {
		return h, nil
	})

	if err != nil {
		b.Fatal("Cannot append a valid route", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := rt.Match("/test/foo")
		if err != nil {
			b.Fatal("Cannot find a valid route;", err)
		}
	}
}
