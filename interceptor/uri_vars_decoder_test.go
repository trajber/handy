package interceptor

import (
	"handy"
	"net"
	"testing"
)

func TestURIVarsBefore(t *testing.T) {
	urivars := map[string]string{
		"s":   "Eita!",
		"b":   "true",
		"i":   "17",
		"i8":  "18",
		"i16": "19",
		"i32": "20",
		"i64": "21",
		"u":   "22",
		"u8":  "23",
		"u16": "24",
		"u32": "25",
		"u64": "26",
		"f32": "27.1",
		"f64": "27.2",
		"ip":  "192.168.0.1",
	}

	handler := struct {
		IntrospectorEmbedded
		handy.DefaultHandler

		S   string  `urivar:"s"`
		B   bool    `urivar:"b"`
		I   int     `urivar:"i"`
		I8  int8    `urivar:"i8"`
		I16 int16   `urivar:"i16"`
		I32 int32   `urivar:"i32"`
		I64 int64   `urivar:"i64"`
		U   uint    `urivar:"u"`
		U8  uint8   `urivar:"u8"`
		U16 uint16  `urivar:"u16"`
		U32 uint32  `urivar:"u32"`
		U64 uint64  `urivar:"u64"`
		F32 float32 `urivar:"f32"`
		F64 float64 `urivar:"f64"`
		IP  net.IP  `urivar:"ip"`
	}{
		DefaultHandler: handy.BuildDefaultHandler(nil, nil, urivars),
	}

	i := NewIntrospector(&handler)
	i.Before()
	u := NewURIVars(urivars, handler.FieldsWithTag("urivar"))
	code := u.Before()

	if code != 0 {
		t.Errorf("Wrong status code. Expecting “0”; found “%d”", code)
	}

	if handler.S != "Eita!" {
		t.Errorf("Wrong value. Expecting “Eita!”; found “%s”", handler.S)
	}

	if handler.B != true {
		t.Errorf("Wrong value. Expecting “true”; found “%t”", handler.B)
	}

	if handler.I != 17 {
		t.Errorf("Wrong value. Expecting “17”; found “%d”", handler.I)
	}

	if handler.I8 != 18 {
		t.Errorf("Wrong value. Expecting “18”; found “%d”", handler.I8)
	}

	if handler.I16 != 19 {
		t.Errorf("Wrong value. Expecting “19”; found “%d”", handler.I16)
	}

	if handler.I32 != 20 {
		t.Errorf("Wrong value. Expecting “20”; found “%d”", handler.I32)
	}

	if handler.I64 != 21 {
		t.Errorf("Wrong value. Expecting “21”; found “%d”", handler.I64)
	}

	if handler.U != 22 {
		t.Errorf("Wrong value. Expecting “22”; found “%d”", handler.U)
	}

	if handler.U8 != 23 {
		t.Errorf("Wrong value. Expecting “23”; found “%d”", handler.U8)
	}

	if handler.U16 != 24 {
		t.Errorf("Wrong value. Expecting “24”; found “%d”", handler.U16)
	}

	if handler.U32 != 25 {
		t.Errorf("Wrong value. Expecting “25”; found “%d”", handler.U32)
	}

	if handler.U64 != 26 {
		t.Errorf("Wrong value. Expecting “26”; found “%d”", handler.U64)
	}

	if handler.F32 != 27.1 {
		t.Errorf("Wrong value. Expecting “27.1”; found “%f”", handler.F32)
	}

	if handler.F64 != 27.2 {
		t.Errorf("Wrong value. Expecting “27.2”; found “%f”", handler.F64)
	}

	if handler.IP == nil || !handler.IP.Equal(net.ParseIP("192.168.0.1")) {
		t.Errorf("Wrong value. Expecting “192.168.0.1”; found “%s”", handler.IP)
	}
}
