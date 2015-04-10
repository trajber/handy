package interceptor

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var tagFormat = regexp.MustCompile(`(\w+):"([^"]+)"`)

type structFields map[string]map[string]reflect.Value

type setFielder interface {
	SetFields(structFields)
}

type Introspector struct {
	NopInterceptor

	structure setFielder
}

func NewIntrospector(st setFielder) *Introspector {
	return &Introspector{structure: st}
}

func (i *Introspector) Before() int {
	st := reflect.ValueOf(i.structure).Elem()
	typ := st.Type()
	fields := make(structFields)

	for i := 0; i < st.NumField(); i++ {
		field := typ.Field(i)

		if field.Tag == "" {
			continue
		}

		for _, ts := range strings.Split(string(field.Tag), " ") {
			tags := tagFormat.FindAllStringSubmatch(ts, -1)

			for _, tagParts := range tags {
				if len(tagParts) != 3 {
					return http.StatusInternalServerError
				}

				name, values := tagParts[1], tagParts[2]

				for _, value := range strings.Split(values, ",") {
					if _, ok := fields[name]; !ok {
						fields[name] = make(map[string]reflect.Value)
					}

					fields[name][value] = st.Field(i)
				}
			}
		}
	}

	i.structure.SetFields(fields)
	return 0
}

type IntrospectorEmbedded struct {
	fields structFields
}

func (i *IntrospectorEmbedded) SetFields(fields structFields) {
	i.fields = fields
}

func (i *IntrospectorEmbedded) Field(tag, value string) interface{} {
	values, found := i.fields[tag]

	if !found {
		return nil
	}

	v, found := values[value]

	if !found || !v.IsValid() || !v.CanInterface() {
		return nil
	}

	if v.Kind() == reflect.Ptr || !v.CanAddr() {
		return v.Interface()
	}

	return v.Addr().Interface()
}

func (i *IntrospectorEmbedded) FieldValues(tag string) (map[string]reflect.Value, bool) {
	values, ok := i.fields[tag]
	return values, ok
}
