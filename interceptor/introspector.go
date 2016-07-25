package interceptor

import (
	"reflect"
	"regexp"
	"strings"
)

var tagFormat = regexp.MustCompile(`(\w+):"([^"]+)"`)

type StructFields map[string]map[string]reflect.Value

type setFielder interface {
	SetFields(StructFields)
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
	fields := make(StructFields)

	i.parse(st, fields)
	i.structure.SetFields(fields)
	return 0
}

func (i *Introspector) parse(st reflect.Value, fields StructFields) {
	typ := st.Type()

	for j := 0; j < st.NumField(); j++ {
		field := typ.Field(j)

		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			i.parse(st.Field(j), fields)
			continue
		}

		if field.Tag == "" {
			continue
		}

		for _, ts := range strings.Split(string(field.Tag), " ") {
			tags := tagFormat.FindAllStringSubmatch(ts, -1)

			for _, tagParts := range tags {
				name, values := tagParts[1], tagParts[2]

				for _, value := range strings.Split(values, ",") {
					if _, ok := fields[name]; !ok {
						fields[name] = make(map[string]reflect.Value)
					}

					fields[name][value] = st.Field(j)
				}
			}
		}
	}
}

func emptyInterface(v reflect.Value) interface{} {
	if !v.IsValid() || !v.CanInterface() {
		return nil
	}

	if v.Kind() == reflect.Ptr || !v.CanAddr() {
		if v.IsNil() {
			return nil
		}

		return v.Interface()
	}

	return v.Addr().Interface()
}

type IntrospectorCompliant struct {
	fields StructFields
}

func (i *IntrospectorCompliant) SetFields(fields StructFields) {
	i.fields = fields
}

func (i *IntrospectorCompliant) SetField(tag, value string, data interface{}) {
	values, found := i.fields[tag]

	if !found {
		return
	}

	f, found := values[value]

	if !found {
		return
	}

	if f.CanSet() {
		f.Set(reflect.ValueOf(data))
	}
}

func (i *IntrospectorCompliant) Field(tag, value string) interface{} {
	values, found := i.fields[tag]

	if !found {
		return nil
	}

	f, found := values[value]

	if !found {
		return nil
	}

	return emptyInterface(f)
}

func (i *IntrospectorCompliant) KeysWithTag(tag string) []string {
	keys := make([]string, 0, len(i.fields[tag]))

	for k := range i.fields[tag] {
		keys = append(keys, k)
	}

	return keys
}
