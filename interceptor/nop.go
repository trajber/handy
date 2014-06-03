package interceptor

import (
	"net/http"
)

type NoBeforeInterceptor struct{}

func (i *NoBeforeInterceptor) Before(w http.ResponseWriter, r *http.Request) {}

type NoAfterInterceptor struct{}

func (i *NoAfterInterceptor) After(w http.ResponseWriter, r *http.Request) {}

type NopInterceptor struct{}

func (i *NopInterceptor) Before(w http.ResponseWriter, r *http.Request) {}
func (i *NopInterceptor) After(w http.ResponseWriter, r *http.Request)  {}
