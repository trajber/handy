package handy

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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

	m := strings.ToLower(r.Method)
	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("request")
		if value == "all" || strings.Contains(value, m) {
			unmarshalURIParams(c.URIParams, st.Field(i))
		}
	}
}

type JSONCodec struct {
	errPosition int
	reqPosition int
	resPosition int
}

func (c *JSONCodec) Encode(w http.ResponseWriter, r *http.Request, h Handler) {
	st := reflect.ValueOf(h).Elem()

	errIface := st.Field(c.errPosition).Interface()
	if c.errPosition >= 0 && errIface != nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(errIface)
		return
	}

	if c.resPosition >= 0 {
		encoder := json.NewEncoder(w)
		encoder.Encode(st.Field(c.resPosition).Interface())
	}
}

func (c *JSONCodec) Decode(w http.ResponseWriter, r *http.Request, h Handler) {
	m := strings.ToLower(r.Method)
	c.parse(m, h)

	if c.reqPosition >= 0 {
		st := reflect.ValueOf(h).Elem()
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(st.Field(c.reqPosition).Addr().Interface())
	}
}

func (c *JSONCodec) parse(m string, h Handler) {
	c.errPosition, c.reqPosition, c.resPosition = -1, -1, -1

	st := reflect.ValueOf(h).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		if field.Tag == "error" {
			c.errPosition = i
			continue
		}

		value := field.Tag.Get("request")
		if value == "all" || strings.Contains(value, m) {
			c.reqPosition = i
			continue
		}

		value = field.Tag.Get("response")
		if value == "all" || strings.Contains(value, m) {
			c.resPosition = i
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
