package handy

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type Codec interface {
	Encode(ctx *Context, h Handler)
	Decode(ctx *Context, h Handler)
}

type NoOpCodec struct{}

func (c *NoOpCodec) Encode(ctx *Context, h Handler) {}
func (c *NoOpCodec) Decode(ctx *Context, h Handler) {}

type ParamCodec struct{}

func (c *ParamCodec) Decode(ctx *Context, h Handler) {
	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("param")

		if value == "" {
			continue
		}

		v := ctx.GetVar(value)
		if v == "" {
			continue
		}

		s := st.FieldByName(field.Name)
		if s.IsValid() && s.CanSet() {
			switch field.Type.Kind() {
			case reflect.String:
				s.SetString(v)
			case reflect.Int:
				i, _ := strconv.ParseInt(v, 10, 0)
				s.SetInt(i)
			}
		}
	}
}

type JSONCodec struct {
	ParamCodec
}

func (c *JSONCodec) Encode(ctx *Context, h Handler) {
	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("codec")
		if value == "response" {
			encoder := json.NewEncoder(ctx.ResponseWriter)
			encoder.Encode(st.Field(i).Interface())
		}
	}
}

func (c *JSONCodec) Decode(ctx *Context, h Handler) {
	c.ParamCodec.Decode(ctx, h)
	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("codec")
		if value == "request" {
			decoder := json.NewDecoder(ctx.Request.Body)
			decoder.Decode(st.Field(i).Interface())
		}
	}
}
