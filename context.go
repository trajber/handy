package handy

import (
	"net/http"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	vars           map[string]string
}

func newContext() *Context {
	ctx := new(Context)
	return ctx
}

func (ctx *Context) GetVar(name string) string {
	return ctx.vars[name]
}
