package handy

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	vars           map[string]string
	Attachment     interface{}
}

func newContext() *Context {
	ctx := new(Context)
	return ctx
}

func (ctx *Context) GetVar(name string) string {
	return ctx.vars[name]
}

func (ctx *Context) Unmarshal(v interface{}) error {
	decoder := json.NewDecoder(ctx.Request.Body)
	return decoder.Decode(v)
}

func (ctx *Context) Marshal(v interface{}) error {
	encoder := json.NewEncoder(ctx.ResponseWriter)
	return encoder.Encode(v)
}
