// Package interceptor provides some ready-to-use interceptors.
package interceptor

import (
	"handy"
	"reflect"
	"regexp"
	"strings"
)

var tagFormat = regexp.MustCompile(`(\w+):"([^"]+)"`)

// Introspector is used by other interceptors for querying fields of structs.
//
// Many interceptors, like QueryString and JSONCodec, need to query tags inside
// the handler using introspecting facilities of the Go language. Introspector
// provides a unifying API to perform such queries.
type Introspector interface {
	handy.Interceptor
	IntrospectorAPI
}

// IntrospectorAPI is the API provided by Introspector to be used by other
// interceptors.
type IntrospectorAPI interface {
	// SetField sets a structure's field tagged in the format `tag:"value"`
	// with the value in the data argument.
	SetField(tag, value string, data interface{})

	// Field queries the value of a structure's field tagged in the format
	// `tag:"value"`
	Field(tag, value string) interface{}

	// KeysWithTag queries the names of the values of a structure's field. For
	// instance, if the field is tagged as `response:"put,post"`,
	// KeysWithTag("response") returns ["put", "post"].
	KeysWithTag(tag string) []string
}

type introspector struct {
	handy.BaseInterceptor

	fields structFields
}

// NewIntrospector creates an Introspector that uses reflection to inspect the
// the structure passed as the second argument. The created Introspector will
// be run by Handy just after the previous interceptor (passed as the first
// argument) executed successfully.
func NewIntrospector(previous handy.Interceptor, structure interface{}) Introspector {
	intro := &introspector{fields: make(structFields)}
	intro.SetPrevious(previous)
	st := reflect.ValueOf(structure).Elem()
	parse(st, intro.fields)

	return intro
}

type structFields map[string]map[string]reflect.Value

func (i *introspector) SetField(tag, value string, data interface{}) {
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

func (i *introspector) Field(tag, value string) interface{} {
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

func (i *introspector) KeysWithTag(tag string) []string {
	keys := make([]string, 0, len(i.fields[tag]))

	for k := range i.fields[tag] {
		keys = append(keys, k)
	}

	return keys
}

func parse(st reflect.Value, fields structFields) {
	typ := st.Type()

	for j := 0; j < st.NumField(); j++ {
		field := typ.Field(j)

		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			parse(st.Field(j), fields)
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
