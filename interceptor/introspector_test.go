package interceptor

import (
	"br/tests"
	"testing"
)

func TestIntrospectorBefore(t *testing.T) {
	type introspectable struct {
		IntrospectorCompliant

		First  int    `one:"um"`
		Second string `two:"dos,dois"`
		Third  int
		Fourth *bool `other:"unrelated" four:"cinco,quatro"`
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

	} else if values := object.KeysWithTag("one"); len(values) != 1 {
		t.Errorf("Wrong number of values for tag “one”: “%d”", len(values))
	}

	if second, ok := object.Field("two", "dois").(*string); !ok || &object.Second != second {
		t.Errorf("It didn't retrieved the right value. Expecting “%#v”; found “%#v”", second, object.Field("second", "dois"))

	} else if values := object.KeysWithTag("two"); len(values) != 2 {
		t.Errorf("Wrong number of values for tag “two”: “%d”", len(values))
	}

	if fourth, ok := object.Field("four", "quatro").(*bool); !ok || object.Fourth != fourth {
		t.Errorf("It didn't retrieved the right value. Expecting “%#v”; found “%#v”", fourth, object.Field("fourth", "quatro"))

	} else if values := object.KeysWithTag("four"); len(values) != 2 {
		t.Errorf("Wrong number of values for tag “four”: “%d”", len(values))
	}

	newFourth := false
	object.Fourth = &newFourth

	if fourth, ok := object.Field("four", "quatro").(*bool); !ok || object.Fourth != fourth {
		t.Errorf("It didn't retrieved the right value. Expecting “%#v”; found “%#v”", fourth, object.Field("fourth", "quatro"))
	}

	// See if the check for a nil value works as expected

	object = introspectable{}
	i = NewIntrospector(&object)
	code = i.Before()

	if code != 0 {
		t.Errorf("Wrong status code. Expecting “0”; found “%d”", code)
	}

	if field := object.Field("four", "quatro"); field != nil {
		t.Errorf("This value is supposed to be nil: %#v", field)
	}

	object.Fourth = &isTrue

	if field := object.Field("four", "quatro"); field == nil {
		t.Errorf("This value isn't supposed to be nil: %#v", field)
	}
}

func TestIntrospectorBeforeCanNotInterface(t *testing.T) {
	object := struct {
		IntrospectorCompliant
		f int `field:"f"`
	}{}
	i := NewIntrospector(&object)
	i.Before()

	if f := object.Field("field", "f"); f != nil {
		t.Errorf("The value %#v is supposed to be nil", f)
	}
}

func TestIntrospectorBeforeUnknownField(t *testing.T) {
	object := struct {
		IntrospectorCompliant
		F int `field:"f"`
	}{}
	i := NewIntrospector(&object)
	i.Before()

	// It shouldn't change the state of the object

	copied := object
	object.SetField("missing", "field", 17)

	if copied.F != object.F {
		t.Errorf("Both objects are expected to be equal:\n%s", tests.Diff(copied, object))
	}

	object.SetField("field", "g", 17)

	if copied.F != object.F {
		t.Errorf("Both objects are expected to be equal:\n%s", tests.Diff(copied, object))
	}

	f := object.Field("missing", "field")

	if f != nil {
		t.Errorf("The value %#v is supposed to be nil", f)
	}
}
