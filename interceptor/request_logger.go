package interceptor

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestLogger struct {
	NoAfterInterceptor

	logger  *log.Logger
	request *http.Request
}

func NewRequestLogger(lg *log.Logger, r *http.Request) *RequestLogger {
	return &RequestLogger{logger: lg}
}

func (l *RequestLogger) Before() int {
	if l.request.Body == nil || l.logger == nil {
		return 0
	}

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, l.request.Body)
	l.request.Body.Close()
	l.request.Body = ioutil.NopCloser(buf)

	return 0
}
