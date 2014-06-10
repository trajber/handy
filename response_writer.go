package handy

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	status   int
	modified bool
	written  bool
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	if !w.Modified() {
		// the first call to Write
		// will trigger an implicit WriteHeader(http.StatusOK).
		w.WriteHeader(http.StatusOK)
	}

	w.ResponseWriter.WriteHeader(w.status)
	w.written = true
	return w.ResponseWriter.Write(b)
}

func (w *ResponseWriter) Written() bool {
	return w.written
}

func (w *ResponseWriter) WriteHeader(s int) {
	w.modified = true
	w.status = s
}

func (w *ResponseWriter) Modified() bool {
	return w.modified
}
