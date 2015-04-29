package interceptor

type NoBeforeInterceptor struct{}

func (i *NoBeforeInterceptor) Before() int {
	return 0
}

type NoAfterInterceptor struct{}

func (i *NoAfterInterceptor) After(status int) int {
	return status
}

type NopInterceptor struct{}

func (i *NopInterceptor) Before() int {
	return 0
}

func (i *NopInterceptor) After(status int) int {
	return status
}

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
