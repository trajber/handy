package handy

import (
	"bytes"
	"net/http"
)

type BufferedResponseWriter struct {
	wire        http.ResponseWriter
	wroteHeader bool
	wroteBody   bool
	flushed     bool
	status      int
	Body        *bytes.Buffer
}

func NewBufferedResponseWriter() *BufferedResponseWriter {
	return &BufferedResponseWriter{
		Body: new(bytes.Buffer),
	}
}

// Header returns the response headers.
func (rw *BufferedResponseWriter) Header() http.Header {
	return rw.wire.Header()
}

// Write always succeeds and writes to rw.Body, if not nil.
func (rw *BufferedResponseWriter) Write(buf []byte) (int, error) {
	if !rw.wroteHeader {
		// note: the first call to Write will trigger an
		// implicit WriteHeader(http.StatusOK).
		rw.WriteHeader(http.StatusOK)
	}

	if rw.Body != nil {
		rw.Body.Write(buf)
		rw.wroteBody = true
	}

	return len(buf), nil
}

func (rw *BufferedResponseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
	}
}

func (rw *BufferedResponseWriter) Flush() {
	if !rw.flushed {
		if !rw.wroteHeader {
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.wire.WriteHeader(rw.status)
		}
	}

	if rw.Body != nil {
		rw.wire.Write(rw.Body.Bytes())
		rw.Body.Reset()
	}

	rw.flushed = true
}
