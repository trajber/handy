package interceptor

import "testing"

func TestIntrospectorBefore(t *testing.T) {
	type introspectable struct {
		IntrospectorEmbedded

		First  int    `one:"um"`
		Second string `two:"dos,dois"`
		Third  int
		Fourth *bool `other: "unrelated" four:"cinco,quatro"`
	}

	isTrue := true
	object := introspectable{
		First:  1,
		Second: "segundo",
		Third:  3,
		Fourth: &isTrue,
	}

	i := NewIntrospector(&object)
	code := i.Before()

	if code != 0 {
		t.Errorf("Wrong status code. Expecting “0”; found “%d”", code)
	}

	if first, ok := object.Field("one", "um").(*int); !ok || &object.First != first {
		t.Errorf("It didn't retrieved the right value. Expecting “%#v”; found “%#v”", first, object.Field("one", "um"))

	} else if values, _ := object.FieldValues("one"); len(values) != 1 {
		t.Errorf("Wrong number of values for tag “one”: “%d”", len(values))
	}

	if second, ok := object.Field("two", "dois").(*string); !ok || &object.Second != second {
		t.Errorf("It didn't retrieved the right value. Expecting “%#v”; found “%#v”", second, object.Field("second", "dois"))

	} else if values, _ := object.FieldValues("two"); len(values) != 2 {
		t.Errorf("Wrong number of values for tag “two”: “%d”", len(values))
	}

	if fourth, ok := object.Field("four", "quatro").(*bool); !ok || object.Fourth != fourth {
		t.Errorf("It didn't retrieved the right value. Expecting “%#v”; found “%#v”", fourth, object.Field("fourth", "quatro"))

	} else if values, _ := object.FieldValues("four"); len(values) != 2 {
		t.Errorf("Wrong number of values for tag “four”: “%d”", len(values))
	}
}
