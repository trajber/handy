package interceptor

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func setValue(ptr interface{}, value string) error {
	switch f := ptr.(type) {
	case nil:
		return nil

	case *string:
		*f = value

	case *bool:
		lower := strings.ToLower(value)
		*f = lower == "true"

	case *int, *int8, *int16, *int32, *int64:
		return setValueInt(ptr, value)

	case *uint, *uint8, *uint16, *uint32, *uint64:
		return setValueUint(ptr, value)

	case *float32, *float64:
		return setValueFloat(ptr, value)

	default:
		return setValueUnmarshaler(ptr, value)
	}

	return nil
}

func setValueInt(ptr interface{}, value string) error {
	n, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return err
	}

	v := reflect.ValueOf(ptr)
	v.Elem().SetInt(n)
	return nil
}

func setValueUint(ptr interface{}, value string) error {
	n, err := strconv.ParseUint(value, 10, 64)

	if err != nil {
		return err
	}

	v := reflect.ValueOf(ptr)
	v.Elem().SetUint(n)
	return nil
}

func setValueFloat(ptr interface{}, value string) error {
	n, err := strconv.ParseFloat(value, 64)

	if err != nil {
		return err
	}

	v := reflect.ValueOf(ptr)
	v.Elem().SetFloat(n)
	return nil
}

func setValueUnmarshaler(ptr interface{}, value string) error {
	u, ok := ptr.(encoding.TextUnmarshaler)
	if !ok {
		return fmt.Errorf("Unsuported value type: %#v", ptr)
	}

	return u.UnmarshalText([]byte(value))
}
