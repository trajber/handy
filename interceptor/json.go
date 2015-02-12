package interceptor

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/trajber/handy"
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
		for {
			if err := decoder.Decode(st.Field(c.reqPosition).Addr().Interface()); err != nil {
				break
			}
		}
	}
}

func (c *JSONCodec) After(w http.ResponseWriter, r *http.Request) {
	st := reflect.ValueOf(c.structure).Elem()

	if c.errPosition >= 0 {
		if elem := st.Field(c.errPosition).Interface(); elem != nil {
			elemType := reflect.TypeOf(elem)
			if elemType.Kind() == reflect.Ptr && !st.Field(c.errPosition).IsNil() {
				encoder := json.NewEncoder(w)
				if err := encoder.Encode(elem); err != nil && handy.ErrorFunc != nil {
					handy.ErrorFunc(err)
				}
				return
			}
		}
	}

	if c.resPosition >= 0 {
		elem := st.Field(c.resPosition).Interface()
		elemType := reflect.TypeOf(elem)
		if elemType.Kind() == reflect.Ptr && st.Field(c.resPosition).IsNil() {
			return
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(elem); err != nil && handy.ErrorFunc != nil {
			handy.ErrorFunc(err)
		}
	}
}

func (c *JSONCodec) parse(m string) {
	c.errPosition, c.reqPosition, c.resPosition = -1, -1, -1

	st := reflect.ValueOf(c.structure).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)

		value := field.Tag.Get("request")
		if value == "all" || strings.Contains(value, m) {
			c.reqPosition = i
			continue
		}

		value = field.Tag.Get("response")

		if value == "all" || strings.Contains(value, m) {
			c.resPosition = i
		} else if value == "error" {
			c.errPosition = i
		}
	}
}
