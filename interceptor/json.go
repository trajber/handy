package interceptor

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

type JSONCodec struct {
	structure interface{}
	err       reflect.Value
	request   reflect.Value
	response  reflect.Value
}

func NewJSONCodec(st interface{}) *JSONCodec {
	return &JSONCodec{structure: st}
}

func (c *JSONCodec) Before(w http.ResponseWriter, r *http.Request) {
	m := strings.ToLower(r.Method)
	c.parse(m)

	if c.request.IsValid() {
		decoder := json.NewDecoder(r.Body)
		for {
			if err := decoder.Decode(c.request.Addr().Interface()); err != nil {
				break
			}
		}
	}
}

func (c *JSONCodec) After(w http.ResponseWriter, r *http.Request) {
	if c.err.IsValid() {
		if elem := c.err.Interface(); elem != nil {
			elemType := reflect.TypeOf(elem)
			if elemType.Kind() == reflect.Ptr && !c.err.IsNil() {
				encoder := json.NewEncoder(w)
				encoder.Encode(elem)
				return
			}
		}
	}

	if c.response.IsValid() {
		elem := c.response.Interface()
		elemType := reflect.TypeOf(elem)
		if elemType.Kind() == reflect.Ptr && c.response.IsNil() {
			return
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(elem)
	}
}

func (c *JSONCodec) parse(m string) {
	st := reflect.ValueOf(c.structure).Elem()

	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := st.Field(i)

		tag := field.Tag.Get("request")
		if tag == "all" || strings.Contains(tag, m) {
			c.request = value
			continue
		}

		tag = field.Tag.Get("response")
		if tag == "all" || strings.Contains(tag, m) {
			c.response = value
		} else if tag == "error" {
			c.err = value
		}
	}
}
