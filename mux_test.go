package handy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInterceptorOrder(t *testing.T) {
	data := []struct {
		description                    string
		interceptors                   InterceptorChain
		shouldBreakAtInterceptorNumber int
	}{
		{
			description: "It should execute all interceptors and the handler",
			interceptors: InterceptorChain{
				&mockInterceptor{},
				&mockInterceptor{},
				&mockInterceptor{},
				&mockInterceptor{},
			},
			shouldBreakAtInterceptorNumber: 1 << 10, // Shouldn't break at all
		},
		{
			description: "It should break at the middle of the chain",
			interceptors: InterceptorChain{
				&mockInterceptor{},
				&mockInterceptor{},
				&brokenBeforeInterceptor{},
				&mockInterceptor{},
			},
			shouldBreakAtInterceptorNumber: 2,
		},
	}

	mux := NewHandy()

	for i, item := range data {
		handleFuncCalled := false
		handler := &mockHandler{
			handleFunc: func() int {
				handleFuncCalled = true
				return http.StatusOK
			},
			interceptors: item.interceptors,
		}

		uri := fmt.Sprintf("/uri/%d", i)
		mux.Handle(uri, func() Handler {
			return handler
		})

		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", uri, nil)

		if err != nil {
			t.Error(err)
		}

		mux.ServeHTTP(w, r)

		for k, interceptor := range item.interceptors {
			interc := interceptor.(MockInterceptor)

			if k <= item.shouldBreakAtInterceptorNumber {
				if !interc.BeforeMethodCalled() {
					t.Errorf("Item %d, “%s”, not calling Before method for interceptor number %d", i, item.description, k)
				}

				if !interc.AfterMethodCalled() {
					t.Errorf("Item %d, “%s”, not calling After method for interceptor number %d", i, item.description, k)
				}

			} else {
				if interc.BeforeMethodCalled() {
					t.Errorf("Item %d, “%s”, calling Before method for interceptor number %d", i, item.description, k)
				}

				if interc.AfterMethodCalled() {
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

type mockInterceptor struct {
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
	Interceptor
	BeforeMethodCalled() bool
	AfterMethodCalled() bool
}

type brokenBeforeInterceptor struct {
	mockInterceptor
}

func (b *brokenBeforeInterceptor) Before() int {
	b.beforeMethodCalled = true
	return http.StatusInternalServerError
}

type brokenAfterInterceptor struct {
	mockInterceptor
}

func (b *brokenAfterInterceptor) After(int) int {
	b.afterMethodCalled = true
	return http.StatusInternalServerError
}

type mockHandler struct {
	DefaultHandler
	handleFunc   func() int
	interceptors InterceptorChain
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

func (m *mockHandler) Interceptors() InterceptorChain {
	return m.interceptors
}
