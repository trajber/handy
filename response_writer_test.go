package handy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	rw = NewBufferedResponseWriter(httptest.NewRecorder())
)

func TestResponseWriterCreation(t *testing.T) {
	if rw.Body == nil {
		t.Fatal("body should not be nil")
	}

	if rw.wire == nil {
		t.Fatal("wire should not be nil")
	}

	if rw.Header() == nil {
		t.Fatal("header should not be nil")
	}
}

func TestResponseWriterModification(t *testing.T) {
	rw.Header().Add("X-Test", "X")
	if rw.Header().Get("X-Test") != "X" {
		t.Fatal("invalid header value")
	}

	data := []byte("test")
	i, err := rw.Write(data)
	if i != len(data) || err != nil {
		t.Fatal("error writing data on the writer")
	}

	if !bytes.Equal(rw.Body.Bytes(), data) {
		t.Fatal("invalid body data")
	}

	rw.WriteHeader(http.StatusTeapot)
	if rw.Status() != http.StatusTeapot {
		t.Fatal("invalid http status code")
	}
}

func TestResponseWriterFlushing(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewBufferedResponseWriter(rec)
	data := []byte("test")
	rw.Write(data)
	rw.Flush()
	if rw.Status() != http.StatusOK {
		t.Fatal("first call to flush should implicitly set status 'ok'")
	}

	if !bytes.Equal(rec.Body.Bytes(), data) {
		t.Fatal("data not sent on the wire")
	}

	if rw.Body.Len() != 0 {
		t.Fatal("body not reseted after flush")
	}
}
