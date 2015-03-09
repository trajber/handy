package handy

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var ParamErrorFunc func(w http.ResponseWriter, r *http.Request, e error)

type paramDecoder struct {
	handler   Handler
	uriParams map[string]string
}

func newParamDecoder(h Handler, uriParams map[string]string) paramDecoder {
	return paramDecoder{handler: h, uriParams: uriParams}
}

func (c *paramDecoder) Decode(w http.ResponseWriter, r *http.Request) {
	st := reflect.ValueOf(c.handler).Elem()
	c.unmarshalURIParams(st, w, r)

	m := strings.ToLower(r.Method)
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		value := field.Tag.Get("request")
		if value == "all" || strings.Contains(value, m) {
			c.unmarshalURIParams(st.Field(i), w, r)
		}
	}
}

func (c *paramDecoder) unmarshalURIParams(st reflect.Value,
	w http.ResponseWriter, r *http.Request) {

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

			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				n, err := strconv.ParseInt(param, 10, 64)
				if err != nil {
					if handleParamError(w, r, err) {
						return
					}
				}
				s.SetInt(n)

			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				n, err := strconv.ParseUint(param, 10, 64)
				if err != nil {
					if handleParamError(w, r, err) {
						return
					}
				}
				s.SetUint(n)

			case reflect.Float32, reflect.Float64:
				n, err := strconv.ParseFloat(param, 64)
				if err != nil {
					if handleParamError(w, r, err) {
						return
					}
				}
				s.SetFloat(n)

			}
		}
	}

	return
}

// returns 'true' when something was written
func handleParamError(w http.ResponseWriter, r *http.Request, err error) bool {
	if ErrorFunc != nil {
		ErrorFunc(err)
	}

	if ParamErrorFunc != nil {
		ParamErrorFunc(w, r, err)
		if bw, ok := w.(*BufferedResponseWriter); ok {
			return bw.somethingWasWritten()
		}
	}

	return false
}
