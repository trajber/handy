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

type BeforeInterceptorFunc func() int

func (i BeforeInterceptorFunc) Before() int {
	return i()
}

func (i BeforeInterceptorFunc) After(int) int {
	return 0
}

type AfterInterceptorFunc func(int) int

func (i AfterInterceptorFunc) Before() int {
	return 0
}

func (i AfterInterceptorFunc) After(status int) int {
	return i(status)
}
