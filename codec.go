package handy

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
)

type Codec interface {
	Encode(http.ResponseWriter, *http.Request, Handler)
	Decode(http.ResponseWriter, *http.Request, Handler)
}

type NoOpCodec struct{}

func (c *NoOpCodec) Encode(http.ResponseWriter, *http.Request, Handler) {}
func (c *NoOpCodec) Decode(http.ResponseWriter, *http.Request, Handler) {}

type ParamCodec struct {
	URIParams map[string]string
}

func (c *ParamCodec) Decode(w http.ResponseWriter, r *http.Request, h Handler) {
	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("param")

		if value == "" {
			continue
		}

		v, ok := c.URIParams[value]
		if !ok {
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

type JSONCodec struct{}

func (c *JSONCodec) Encode(w http.ResponseWriter, r *http.Request, h Handler) {
	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("codec")
		if value == "response" {
			encoder := json.NewEncoder(w)
			encoder.Encode(st.Field(i).Interface())
		}
	}
}

func (c *JSONCodec) Decode(w http.ResponseWriter, r *http.Request, h Handler) {
	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("codec")
		if value == "request" {
			decoder := json.NewDecoder(r.Body)
			decoder.Decode(st.Field(i).Interface())
		}
	}
}
