package handy

import (
	"bytes"
	"net/http"
)

type BufferedResponseWriter struct {
	flushed bool
	status  int
	wire    http.ResponseWriter
	Body    *bytes.Buffer
}

func NewBufferedResponseWriter(w http.ResponseWriter) *BufferedResponseWriter {
	return &BufferedResponseWriter{
		Body: new(bytes.Buffer),
		wire: w,
	}
}

// Header returns the response headers.
func (rw *BufferedResponseWriter) Header() http.Header {
	return rw.wire.Header()
}

// Write always succeeds and writes to rw.Body, if not nil.
func (rw *BufferedResponseWriter) Write(buf []byte) (int, error) {
	return rw.Body.Write(buf)
}

func (rw *BufferedResponseWriter) Status() int {
	return rw.status
}

func (rw *BufferedResponseWriter) WriteHeader(code int) {
	rw.status = code
}

func (rw *BufferedResponseWriter) Flush() {
	if rw.status == 0 {
		rw.WriteHeader(http.StatusOK)
	}

	if !rw.flushed {
		rw.wire.WriteHeader(rw.status)
	}

	rw.wire.Write(rw.Body.Bytes())
	rw.Body.Reset()

	rw.flushed = true
}
