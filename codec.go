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

type NopCodec struct{}

func (c *NopCodec) Encode(http.ResponseWriter, *http.Request, Handler) {}
func (c *NopCodec) Decode(http.ResponseWriter, *http.Request, Handler) {}

type ParamCodec struct {
	URIParams map[string]string
}

func (c *ParamCodec) Decode(w http.ResponseWriter, r *http.Request, h Handler) {
	unmarshalURIParams(c.URIParams, reflect.ValueOf(h).Elem())

	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("codec")
		if value == "request" {
			unmarshalURIParams(c.URIParams, st.Field(i))
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
		if value == "request" && r.Body != nil {
			decoder := json.NewDecoder(r.Body)
			decoder.Decode(st.Field(i).Addr().Interface())
		}
	}
}

func unmarshalURIParams(uriParams map[string]string, st reflect.Value) {
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("param")

		if value == "" {
			continue
		}

		param, ok := uriParams[value]
		if !ok {
			continue
		}

		s := st.FieldByName(field.Name)
		if s.IsValid() && s.CanSet() {
			switch field.Type.Kind() {
			case reflect.String:
				s.SetString(param)
			case reflect.Int:
				i, err := strconv.ParseInt(param, 10, 0)
				if err != nil {
					Logger.Println(err)
					continue
				}
				s.SetInt(i)
			}
		}
	}
}
