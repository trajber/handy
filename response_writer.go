package handy

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	if !w.Written() {
		// the first call to Write
		// will trigger an implicit WriteHeader(http.StatusOK).
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(b)
}

func (w *ResponseWriter) Written() bool {
	return w.status != 0
}

func (w *ResponseWriter) WriteHeader(s int) {
	w.ResponseWriter.WriteHeader(s)
	w.status = s
}
