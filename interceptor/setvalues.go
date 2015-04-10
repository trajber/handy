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
	case *string:
		*f = value

	case *bool:
		lower := strings.ToLower(value)
		*f = lower == "true"

	case *int, *int8, *int16, *int32, *int64:
		n, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			return err
		}

		v := reflect.ValueOf(ptr)
		v.Elem().SetInt(n)

	case *uint, *uint8, *uint16, *uint32, *uint64:
		n, err := strconv.ParseUint(value, 10, 64)

		if err != nil {
			return err
		}

		v := reflect.ValueOf(ptr)
		v.Elem().SetUint(n)

	case *float32, *float64:
		n, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return err
		}

		v := reflect.ValueOf(ptr)
		v.Elem().SetFloat(n)

	default:
		u, ok := ptr.(encoding.TextUnmarshaler)

		if !ok {
			return fmt.Errorf("Unsuported value type: %#v", ptr)
		}

		return u.UnmarshalText([]byte(value))
	}

	return nil
}
