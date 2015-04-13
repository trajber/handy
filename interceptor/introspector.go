package interceptor

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var tagFormat = regexp.MustCompile(`(\w+):"([^"]+)"`)

type structFields map[string]map[string]interface{}

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
						fields[name] = make(map[string]interface{})
					}

					fields[name][value] = emptyInterface(st.Field(i))
				}
			}
		}
	}

	i.structure.SetFields(fields)
	return 0
}

func emptyInterface(v reflect.Value) interface{} {
	if !v.IsValid() || !v.CanInterface() {
		return nil
	}

	if v.Kind() == reflect.Ptr || !v.CanAddr() {
		return v.Interface()
	}

	return v.Addr().Interface()
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

	f, found := values[value]

	if !found {
		return nil
	}

	return f
}

func (i *IntrospectorEmbedded) FieldsWithTag(tag string) map[string]interface{} {
	return i.fields[tag]
}
