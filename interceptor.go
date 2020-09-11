package handy

type Interceptor interface {
	Before() int
	After(int) int
	SetContext(Context)
	SetPrevious(Interceptor)
	previous() Interceptor
}

// ProtoInterceptor adds support for a struct to be used as an interceptor. It
// is expected to be embedded in all interceptors.
type ProtoInterceptor struct {
	Context

	// Previous is set by the constructor of an interceptor as the interceptor
	// to be executed just before it in the interceptor chain of execution.
	//
	// For instance, if we have the following interceptors:
	//
	//     type A struct {
	//	       ProtoInterceptor
	//     }
	//
	//     func NewA(previous Interceptor) *A {
	//	       return &A{ Previous: previous }
	//     }
	//
	//     type B struct {
	//	       ProtoInterceptor
	//     }
	//
	//     func NewB(previous Interceptor) *B {
	//	       return &B{ Previous: previous }
	//     }
	//
	//     type C struct {
	//	       ProtoInterceptor
	//     }
	//
	//     func NewC(previous Interceptor) *C {
	//	       return &C{ Previous: previous }
	//     }
	//
	// And, for a particular scenario, we want to execute A first,
	// then B, then C, we must write:
	//
	//     func NewHandler() (Handler, Interceptor) {
	//	       return &SomeHandler{}, NewC(NewB(NewA(nil)))
	//     }
	//
	// This is usefull to enforce dependencies at the type level. Suppose
	// that C, if included in the chain of execution, must always run
	// after B has already run. In such a case, we can write:
	//
	//     func NewA(previous Interceptor) *A {
	//	       return &A{ Previous: previous }
	//     }
	//
	//	   func NewB(previous Interceptor) *B {
	//	       return &B{ Previous: previous }
	//     }
	//
	//     func NewC(previous &B) *C {
	//         if previous == nil {
	//             panic("Interceptor C must run after an interceptor B")
	//         }
	//
	//	       return &C{ Previous: previous }
	//     }
	//
	// And the type system garanties we can not run C if we haven't yet run B:
	//
	//     func NewHandler() (Handler, Interceptor) {
	//	       return &SomeHandler{}, NewB(NewC(NewA(nil))) // <-- Doesn't compile!
	//     }
	previousInterceptor Interceptor
}

func (i *ProtoInterceptor) Before() int {
	return 0
}

func (i *ProtoInterceptor) After(status int) int {
	return status
}

func (i *ProtoInterceptor) previous() Interceptor {
	return i.previousInterceptor
}

func (i *ProtoInterceptor) SetPrevious(previous Interceptor) {
	i.previousInterceptor = previous
}

func (i *ProtoInterceptor) SetContext(c Context) {
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
