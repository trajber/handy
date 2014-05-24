package handy

import (
	"reflect"
	"strings"
)

func JSONMarshaller(ctx *Context, h Handler) {
	verb := strings.ToLower(ctx.Request.Method)
	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get(verb)
		if value == "response" {
			ctx.Marshal(st.Field(i).Interface())
		}
	}
}

func JSONUnmarshaller(ctx *Context, h Handler) {
	verb := strings.ToLower(ctx.Request.Method)
	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get(verb)
		if value == "request" {
			ctx.Unmarshal(st.Field(i).Interface())
		}
	}
}
