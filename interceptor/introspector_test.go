package interceptor_test

import (
	"github.com/kylelemons/godebug/pretty"
	"handy"
	"handy/interceptor"
	"testing"
)

func TestIntrospector(t *testing.T) {
	type introspectable struct {
		handy.ProtoHandler
		interceptor.IntrospectorAPI

		First  int    `one:"um"`
		Second string `two:"dos,dois"`
		Third  int
		Fourth *bool `other:"unrelated" four:"cinco,quatro"`
	}

	newIntrospectable := func(object *introspectable) (handy.Handler, handy.Interceptor) {
		intro := interceptor.NewIntrospector(nil, object)
		object.IntrospectorAPI = intro

		return object, intro
	}

	isTrue := true
	object := &introspectable{
		First:  1,
		Second: "segundo",
		Third:  3,
		Fourth: &isTrue,
	}

	newIntrospectable(object)

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

	object = &introspectable{}
	newIntrospectable(object)

	if field := object.Field("four", "quatro"); field != nil {
		t.Errorf("This value is supposed to be nil: %#v", field)
	}

	object.Fourth = &isTrue

	if field := object.Field("four", "quatro"); field == nil {
		t.Errorf("This value isn't supposed to be nil: %#v", field)
	}
}

func TestIntrospectorCanNotInterface(t *testing.T) {
	object := struct {
		interceptor.IntrospectorAPI
		f int `field:"f"`
	}{}
	intro := interceptor.NewIntrospector(nil, &object)
	object.IntrospectorAPI = intro

	if f := object.Field("field", "f"); f != nil {
		t.Errorf("The value %#v is supposed to be nil", f)
	}
}

func TestIntrospectorUnknownField(t *testing.T) {
	object := struct {
		interceptor.IntrospectorAPI
		F int `field:"f"`
	}{}
	intro := interceptor.NewIntrospector(nil, &object)
	object.IntrospectorAPI = intro

	// It shouldn't change the state of the object

	copied := object
	object.SetField("missing", "field", 17)

	if copied.F != object.F {
		t.Errorf("Both objects are expected to be equal:\n%s", pretty.Compare(copied, object))
	}

	object.SetField("field", "g", 17)

	if copied.F != object.F {
		t.Errorf("Both objects are expected to be equal:\n%s", pretty.Compare(copied, object))
	}

	f := object.Field("missing", "field")

	if f != nil {
		t.Errorf("The value %#v is supposed to be nil", f)
	}
}

func TestIntrospectorEmbedded(t *testing.T) {
	object := struct {
		interceptor.IntrospectorAPI
		dummy
	}{}
	intro := interceptor.NewIntrospector(nil, &object)
	object.IntrospectorAPI = intro

	if _, ok := object.Field("field", "f").(*int); !ok {
		t.Error("It didn't identify the object")

	} else if values := object.KeysWithTag("field"); len(values) != 1 {
		t.Errorf("Wrong number of values for tag “field”: “%d”", len(values))
	}
}

type dummy struct {
	F int `field:"f"`
}
