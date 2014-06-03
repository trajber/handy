package interceptor

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type requestLogger struct {
	logger *log.Logger
	NoAfterInterceptor
}

func NewRequestLogger(lg *log.Logger) *requestLogger {
	return &requestLogger{logger: lg}
}

func (l *requestLogger) Before(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil || l.logger == nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(buf)
}
