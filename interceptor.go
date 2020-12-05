package handy

// Interceptor is the way handy implements the decorator pattern.
//
// An interceptor decorates a handler and can execute custom actions before and
// after it handles the request. It can also specifies other interceptors to
// decorate itself, effectively building a nested structure of decorators.
//
// Suppose we have the following interceptors:
//
//     type A struct {
//         BaseInterceptor
//     }
//
//     func NewA(previous Interceptor) *A {
//         a := new(A)
//         a.SetPrevious(previous)
//         return a
//     }
//
//     type B struct {
//         BaseInterceptor
//     }
//
//     func NewB(previous Interceptor) *B {
//         b := new(A)
//         b.SetPrevious(previous)
//         return b
//     }
//
//     type C struct {
//         BaseInterceptor
//     }
//
//     func NewC(previous Interceptor) *C {
//         c := new(A)
//         c.SetPrevious(previous)
//         return c
//     }
//
// And, for a particular scenario, we want to execute A first,
// then B, then C. So we must write:
//
//     func NewHandler() (Handler, Interceptor) {
//         a := NewA(nil)
//         b := NewB(a)
//         c := NewC(b)
//         return new(SomeHandler), c
//     }
//
// In such a setup, when a request arrives say, with a PUT method, the following
// execution chain is performed:
//
//     a.Before, b.Before, c.Before, handler.Put, c.After, b.After, a.After
//
// If any of the interceptors' Before method returns a non-zero value, the
// execution chain is interrupted and neither the subsequent interceptors nor
// the handler are called. It acts as if the following code were executed:
//
//     result := a.Before()
//
//     if result == 0 {
//         result = b.Before()
//
//         if result == 0 {
//             result = c.Before()
//
//             if result == 0 {
//                 result = handler.Put()
//             }
//
//             result = c.After(result)
//         }
//
//         result = b.After(result)
//     }
//
//     a.After(result)
type Interceptor interface {
	Before() int
	After(int) int
	SetPrevious(Interceptor)
	SetContext(Context)
	previous() Interceptor
}

// BaseInterceptor is a prototype implementation for an interceptor. It must be
// embedded in all interceptors to make them compatible with Handy.
type BaseInterceptor struct {
	// Context allows the interceptor to interact with the request and
	// eventually write custom responses.
	Context

	previousInterceptor Interceptor
}

// Before is the first interceptor's method called. A returned value of zero
// signals that everything is OK and the execution chain should continue. A
// non-zero value interrupts the execution chain and its After method, as well
// as all After methods of previous interceptors, are called with such a value
// as an argument.
func (i *BaseInterceptor) Before() int {
	return 0
}

// After is the last interceptor's method called. It receives as argument the
// value propagated by the method calls in the execution chain.
func (i *BaseInterceptor) After(status int) int {
	return status
}

func (i *BaseInterceptor) previous() Interceptor {
	return i.previousInterceptor
}

// SetPrevious registers the provided interceptor as a decorator of the current
// one. The framework, then, will arrange the execution chain as follows:
//
//     previous.Before, i.Before, handler.Method, i.After, previous.After.
func (i *BaseInterceptor) SetPrevious(previous Interceptor) {
	i.previousInterceptor = previous
}

// SetContext is used internally by the framework to set Context information on
// each interceptor. It's not meant to be called by the user, but it's exported
// as a convenience to inject mock data during your tests.
func (i *BaseInterceptor) SetContext(c Context) {
	i.Context = c

	// Recursively set context for all interceptors in the list
	if i.previousInterceptor != nil {
		i.previousInterceptor.SetContext(c)
	}
}

// buildChain makes a new slice of interceptors, in the reverse order of
// the one specified by the user.
func buildChain(interceptor Interceptor) []Interceptor {
	// Pre-allocate some space for performance reasons.
	chain := make([]Interceptor, 0, 8)

	for i := interceptor; i != nil; i = i.previous() {
		chain = append(chain, i)
	}

	return chain
}
