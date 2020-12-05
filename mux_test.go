package handy_test

import (
	"fmt"
	"handy"
	"net/http"
	"net/http/httptest"
	"testing"
)

type interceptorConstructor func(handy.Interceptor) MockInterceptor

func TestInterceptorOrder(t *testing.T) {
	data := []struct {
		description                    string
		interceptors                   []interceptorConstructor
		shouldBreakAtInterceptorNumber int
	}{
		{
			description: "It should execute all interceptors and the handler",
			interceptors: []interceptorConstructor{
				newMockInterceptor,
				newMockInterceptor,
				newMockInterceptor,
				newMockInterceptor,
			},
			shouldBreakAtInterceptorNumber: 1 << 10, // Shouldn't break at all
		},
		{
			description: "It should break at the middle of the chain",
			interceptors: []interceptorConstructor{
				newMockInterceptor,
				newMockInterceptor,
				newBrokenBeforeInterceptor,
				newMockInterceptor,
			},
			shouldBreakAtInterceptorNumber: 2,
		},
	}

	mux := handy.New()

	for i, item := range data {
		handleFuncCalled := false
		handler := &mockHandler{
			handleFunc: func() int {
				handleFuncCalled = true
				return http.StatusOK
			},
		}

		uri := fmt.Sprintf("/uri/%d", i)
		var interceptors []MockInterceptor
		mux.Handle(uri, func() (handy.Handler, handy.Interceptor) {
			var interceptor MockInterceptor

			for _, constructor := range item.interceptors {
				interceptor = constructor(interceptor)
				interceptors = append(interceptors, interceptor)
			}

			return handler, interceptor
		})

		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", uri, nil)

		if err != nil {
			t.Error(err)
		}

		mux.ServeHTTP(w, r)

		for k, interceptor := range interceptors {
			if k <= item.shouldBreakAtInterceptorNumber {
				if !interceptor.BeforeMethodCalled() {
					t.Errorf("Item %d, “%s”, not calling Before method for interceptor number %d", i, item.description, k)
				}

				if !interceptor.AfterMethodCalled() {
					t.Errorf("Item %d, “%s”, not calling After method for interceptor number %d", i, item.description, k)
				}

			} else {
				if interceptor.BeforeMethodCalled() {
					t.Errorf("Item %d, “%s”, calling Before method for interceptor number %d", i, item.description, k)
				}

				if interceptor.AfterMethodCalled() {
					t.Errorf("Item %d, “%s”, calling After method for interceptor number %d", i, item.description, k)
				}
			}
		}

		if len(item.interceptors) < item.shouldBreakAtInterceptorNumber {
			if !handleFuncCalled {
				t.Errorf("Item %d, “%s”, not calling handler", i, item.description)
			}
		} else {
			if handleFuncCalled {
				t.Errorf("Item %d, “%s”, calling handler", i, item.description)
			}
		}
	}
}

func newMockInterceptor(previous handy.Interceptor) MockInterceptor {
	i := new(mockInterceptor)
	i.SetPrevious(previous)
	return i
}

type mockInterceptor struct {
	handy.BaseInterceptor

	beforeMethodCalled bool
	afterMethodCalled  bool
}

func (m *mockInterceptor) Before() int {
	m.beforeMethodCalled = true
	return 0
}

func (m *mockInterceptor) After(int) int {
	m.afterMethodCalled = true
	return 0
}

func (m *mockInterceptor) BeforeMethodCalled() bool {
	return m.beforeMethodCalled
}

func (m *mockInterceptor) AfterMethodCalled() bool {
	return m.afterMethodCalled
}

type MockInterceptor interface {
	handy.Interceptor
	BeforeMethodCalled() bool
	AfterMethodCalled() bool
}

func newBrokenBeforeInterceptor(previous handy.Interceptor) MockInterceptor {
	i := new(brokenBeforeInterceptor)
	i.SetPrevious(previous)
	return i
}

type brokenBeforeInterceptor struct {
	mockInterceptor
}

func (b *brokenBeforeInterceptor) Before() int {
	b.beforeMethodCalled = true
	return http.StatusInternalServerError
}

func newBrokenAfterInterceptor(previous handy.Interceptor) MockInterceptor {
	i := new(brokenAfterInterceptor)
	i.SetPrevious(previous)
	return i
}

type brokenAfterInterceptor struct {
	mockInterceptor
}

func (b *brokenAfterInterceptor) After(int) int {
	b.afterMethodCalled = true
	return http.StatusInternalServerError
}

type mockHandler struct {
	handy.BaseHandler

	handleFunc   func() int
	methodCalled string
}

func (m *mockHandler) Get() int {
	m.methodCalled = "GET"
	return m.handleFunc()
}

func (m *mockHandler) Post() int {
	m.methodCalled = "POST"
	return m.handleFunc()
}

func (m *mockHandler) Put() int {
	m.methodCalled = "PUT"
	return m.handleFunc()
}

func (m *mockHandler) Delete() int {
	m.methodCalled = "DELETE"
	return m.handleFunc()
}

func (m *mockHandler) Patch() int {
	m.methodCalled = "PATCH"
	return m.handleFunc()
}

func (m *mockHandler) Head() int {
	m.methodCalled = "HEAD"
	return m.handleFunc()
}
