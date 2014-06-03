package interceptor

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

type JSONCodec struct {
	structure   interface{}
	errPosition int
	reqPosition int
	resPosition int
}

func NewJSONCodec(st interface{}) *JSONCodec {
	return &JSONCodec{structure: st}
}

func (c *JSONCodec) Before(w http.ResponseWriter, r *http.Request) {
	m := strings.ToLower(r.Method)
	c.parse(m)

	if c.reqPosition >= 0 {
		st := reflect.ValueOf(c.structure).Elem()
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(st.Field(c.reqPosition).Addr().Interface())
	}
}

func (c *JSONCodec) After(w http.ResponseWriter, r *http.Request) {
	st := reflect.ValueOf(c.structure).Elem()

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

func (c *JSONCodec) parse(m string) {
	c.errPosition, c.reqPosition, c.resPosition = -1, -1, -1

	st := reflect.ValueOf(c.structure).Elem()
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
