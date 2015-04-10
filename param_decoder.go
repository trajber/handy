package handy

import (
	"encoding"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type paramDecoder struct {
	handler   Handler
	uriParams map[string]string
}

func newParamDecoder(h Handler, uriParams map[string]string) paramDecoder {
	return paramDecoder{handler: h, uriParams: uriParams}
}

func (c *paramDecoder) Decode(w http.ResponseWriter, r *http.Request) {
	st := reflect.ValueOf(c.handler).Elem()
	c.unmarshalURIParams(st, w)

	m := strings.ToLower(r.Method)
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("request")
		if value == "all" || strings.Contains(value, m) {
			c.unmarshalURIParams(st.Field(i), w)
		}
	}
}

func (c *paramDecoder) unmarshalURIParams(st reflect.Value, w http.ResponseWriter) {
	if st.Kind() == reflect.Ptr {
		return
	}

	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("param")

		if value == "" {
			continue
		}

		param, ok := c.uriParams[value]
		if !ok {
			continue
		}

		s := st.Field(i)
		if s.IsValid() && s.CanSet() {
			switch field.Type.Kind() {
			case reflect.String:
				s.SetString(param)

			case reflect.Bool:
				lower := strings.ToLower(param)
				s.SetBool(lower == "true")

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				n, err := strconv.ParseInt(param, 10, 64)
				if err != nil {
					if ErrorFunc != nil {
						ErrorFunc(err)
					}
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				s.SetInt(n)

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				n, err := strconv.ParseUint(param, 10, 64)
				if err != nil {
					if ErrorFunc != nil {
						ErrorFunc(err)
					}
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				s.SetUint(n)

			case reflect.Float32, reflect.Float64:
				n, err := strconv.ParseFloat(param, 64)
				if err != nil {
					if ErrorFunc != nil {
						ErrorFunc(err)
					}
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				s.SetFloat(n)
			default:
				// check if the structure implements TextUnmarshaler
				if !s.CanAddr() {
					continue
				}

				u, ok := s.Addr().Interface().(encoding.TextUnmarshaler)
				if !ok {
					continue
				}

				if err := u.UnmarshalText([]byte(param)); err != nil {
					if ErrorFunc != nil {
						ErrorFunc(err)
					}
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		}
	}

	return
}
