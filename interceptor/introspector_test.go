package interceptor_test

import (
	"handy"
	"handy/interceptor"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestIntrospector(t *testing.T) {
	type introspectable struct {
		handy.BaseHandler

		First  int    `one:"um"`
		Second string `two:"dos,dois"`
		Third  int
		Fourth *bool `other:"unrelated" four:"cinco,quatro"`
	}

	isTrue := true
	object := &introspectable{
		First:  1,
		Second: "segundo",
		Third:  3,
		Fourth: &isTrue,
	}

	intro := interceptor.NewIntrospector(nil, object)

	if first, ok := intro.Field("one", "um").(*int); !ok || &object.First != first {
		t.Errorf("It didn't retrieved the right value. Expecting “%#v”; found “%#v”", first, intro.Field("one", "um"))

	} else if values := intro.KeysWithTag("one"); len(values) != 1 {
		t.Errorf("Wrong number of values for tag “one”: “%d”", len(values))
	}

	if second, ok := intro.Field("two", "dois").(*string); !ok || &object.Second != second {
		t.Errorf("It didn't retrieved the right value. Expecting “%#v”; found “%#v”", second, intro.Field("second", "dois"))

	} else if values := intro.KeysWithTag("two"); len(values) != 2 {
		t.Errorf("Wrong number of values for tag “two”: “%d”", len(values))
	}

	if fourth, ok := intro.Field("four", "quatro").(*bool); !ok || object.Fourth != fourth {
		t.Errorf("It didn't retrieved the right value. Expecting “%#v”; found “%#v”", fourth, intro.Field("fourth", "quatro"))

	} else if values := intro.KeysWithTag("four"); len(values) != 2 {
		t.Errorf("Wrong number of values for tag “four”: “%d”", len(values))
	}

	newFourth := false
	object.Fourth = &newFourth

	if fourth, ok := intro.Field("four", "quatro").(*bool); !ok || object.Fourth != fourth {
		t.Errorf("It didn't retrieved the right value. Expecting “%#v”; found “%#v”", fourth, intro.Field("fourth", "quatro"))
	}

	// See if the check for a nil value works as expected
	object = new(introspectable)
	intro = interceptor.NewIntrospector(nil, object)

	if field := intro.Field("four", "quatro"); field != nil {
		t.Errorf("This value is supposed to be nil: %#v", field)
	}

	object.Fourth = &isTrue

	if field := intro.Field("four", "quatro"); field == nil {
		t.Errorf("This value isn't supposed to be nil: %#v", field)
	}
}

func TestIntrospectorCanNotInterface(t *testing.T) {
	object := struct {
		f int `field:"f"`
	}{}
	intro := interceptor.NewIntrospector(nil, &object)

	if f := intro.Field("field", "f"); f != nil {
		t.Errorf("The value %#v is supposed to be nil", f)
	}
}

func TestIntrospectorUnknownField(t *testing.T) {
	object := struct {
		F int `field:"f"`
	}{}
	intro := interceptor.NewIntrospector(nil, &object)

	// It shouldn't change the state of the object
	copied := object
	intro.SetField("missing", "field", 17)

	if copied.F != object.F {
		t.Errorf("Both objects are expected to be equal:\n%s", pretty.Compare(copied, object))
	}

	intro.SetField("field", "g", 17)

	if copied.F != object.F {
		t.Errorf("Both objects are expected to be equal:\n%s", pretty.Compare(copied, object))
	}

	f := intro.Field("missing", "field")

	if f != nil {
		t.Errorf("The value %#v is supposed to be nil", f)
	}
}

func TestIntrospectorEmbedded(t *testing.T) {
	object := struct {
		dummy
	}{}
	intro := interceptor.NewIntrospector(nil, &object)

	if _, ok := intro.Field("field", "f").(*int); !ok {
		t.Error("It didn't identify the object")

	} else if values := intro.KeysWithTag("field"); len(values) != 1 {
		t.Errorf("Wrong number of values for tag “field”: “%d”", len(values))
	}
}

type dummy struct {
	F int `field:"f"`
}
