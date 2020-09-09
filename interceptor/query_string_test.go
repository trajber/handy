package interceptor_test

import (
	"errors"
	"fmt"
	"handy"
	"handy/interceptor"
	"net"
	"net/http"
	"net/url"
	"testing"
)

type queryStringHandler struct {
	handy.ProtoHandler
	interceptor.QueryStringAPI

	request *http.Request

	S       string      `query:"s"`
	B       bool        `query:"b"`
	I       int         `query:"i"`
	I8      int8        `query:"i8"`
	I16     int16       `query:"i16"`
	I32     int32       `query:"i32"`
	I64     int64       `query:"i64"`
	Iempty  int         `query:"iempty"`
	U       uint        `query:"u"`
	U8      uint8       `query:"u8"`
	U16     uint16      `query:"u16"`
	U32     uint32      `query:"u32"`
	U64     uint64      `query:"u64"`
	F32     float32     `query:"f32"`
	F64     float64     `query:"f64"`
	IP      net.IP      `query:"ip"`
	Unknown struct{}    `query:"unknown"`
	Custom  *customType `query:"custom"`
}

type customType struct {
	mockUnmarshalText func([]byte) error
}

func (c customType) UnmarshalText(data []byte) error {
	return c.mockUnmarshalText(data)
}

func TestQueryStringBefore(t *testing.T) {
	data := []struct {
		description    string
		queryString    string
		request        *http.Request
		customTypeMock customType
		expected       queryStringHandler
		expectedStatus int
	}{
		{
			description: "it should parse the parameters to the correct types",
			queryString: "s=Eita!&b=true&i=17&i8=18&i16=19&i32=20&i64=21&u=22&u8=23&u16=24&u32=25&u64=26&f32=27.1&f64=27.2&ip=192.168.0.1",
			expected: queryStringHandler{
				S:      "Eita!",
				B:      true,
				I:      17,
				I8:     18,
				I16:    19,
				I32:    20,
				I64:    21,
				Iempty: 0,
				U:      22,
				U8:     23,
				U16:    24,
				U32:    25,
				U64:    26,
				F32:    27.1,
				F64:    27.2,
				IP:     net.ParseIP("192.168.0.1"),
			},
			expectedStatus: 0,
		},
		{
			description:    "it should fail to load an invalid int",
			queryString:    "i=xxxx",
			expectedStatus: http.StatusBadRequest,
		},
		{
			description:    "it should fail to load an invalid uint",
			queryString:    "u=xxxx",
			expectedStatus: http.StatusBadRequest,
		},
		{
			description:    "it should fail to load an invalid float",
			queryString:    "f32=xxxx",
			expectedStatus: http.StatusBadRequest,
		},
		{
			description:    "it should ignore a parameter that does not exist in the handler",
			queryString:    "idontexist=123",
			expectedStatus: 0,
		},
		{
			description: "it should ignore a parameter without value",
			request: &http.Request{
				Form: url.Values(map[string][]string{
					"s":       []string{"Hello!"},
					"novalue": []string{},
				}),
			},
			expected: queryStringHandler{
				S: "Hello!",
			},
			expectedStatus: 0,
		},
		{
			description:    "it should fail to load to an unsupported type",
			queryString:    "unknown=2015-05-07",
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "it should load a custom value correctly",
			queryString: "custom=abc",
			customTypeMock: customType{
				mockUnmarshalText: func(data []byte) error {
					return nil
				},
			},
			expectedStatus: 0,
		},
		{
			description: "it should fail to load a custom value with an invalid error type",
			queryString: "custom=abc",
			customTypeMock: customType{
				mockUnmarshalText: func(data []byte) error {
					return fmt.Errorf("I'm a crazy error!")
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "it should fail to load a custom value",
			queryString: "custom=abc",
			customTypeMock: customType{
				mockUnmarshalText: func(data []byte) error {
					return errors.New("Eta, erro doido!")
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for i, item := range data {
		var request *http.Request
		if item.request == nil {
			var err error
			request, err = http.NewRequest("GET", "http://um.com.br?"+item.queryString, nil)
			if err != nil {
				t.Fatal(err)
			}

		} else {
			request = item.request
		}

		handler := &queryStringHandler{request: request}
		handler.Custom = &item.customTypeMock
		intro := interceptor.NewIntrospector(nil, handler)
		query := interceptor.NewQueryString(intro)
		handler.QueryStringAPI = query

		// The context is set automatically by the framework but on
		// the tests we need to set it manually
		ctx := handy.Context{Request: request}
		handler.SetContext(ctx)
		query.SetContext(ctx)

		status := query.Before()

		if status != item.expectedStatus {
			t.Errorf("Item %d, “%s”: mismatch HTTP status. Expecting “%d”; found “%d”", i, item.description, item.expectedStatus, status)
		}

		if handler.S != item.expected.S {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%s”; found “%s”", i, item.description, item.expected.S, handler.S)
		}

		if handler.B != item.expected.B {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%t”; found “%t”", i, item.description, item.expected.B, handler.B)
		}

		if handler.I != item.expected.I {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%d”; found “%d”", i, item.description, item.expected.I, handler.I)
		}

		if handler.I8 != item.expected.I8 {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%d”; found “%d”", i, item.description, item.expected.I8, handler.I8)
		}

		if handler.I16 != item.expected.I16 {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%d”; found “%d”", i, item.description, item.expected.I16, handler.I16)
		}

		if handler.I32 != item.expected.I32 {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%d”; found “%d”", i, item.description, item.expected.I32, handler.I32)
		}

		if handler.I64 != item.expected.I64 {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%d”; found “%d”", i, item.description, item.expected.I64, handler.I64)
		}

		if handler.U != item.expected.U {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%d”; found “%d”", i, item.description, item.expected.U, handler.U)
		}

		if handler.U8 != item.expected.U8 {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%d”; found “%d”", i, item.description, item.expected.U8, handler.U8)
		}

		if handler.U16 != item.expected.U16 {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%d”; found “%d”", i, item.description, item.expected.U16, handler.U16)
		}

		if handler.U32 != item.expected.U32 {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%d”; found “%d”", i, item.description, item.expected.U32, handler.U32)
		}

		if handler.U64 != item.expected.U64 {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%d”; found “%d”", i, item.description, item.expected.U64, handler.U64)
		}

		if handler.F32 != item.expected.F32 {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%f”; found “%f”", i, item.description, item.expected.F32, handler.F32)
		}

		if handler.F64 != item.expected.F64 {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%f”; found “%f”", i, item.description, item.expected.F64, handler.F64)
		}

		if !handler.IP.Equal(item.expected.IP) {
			t.Errorf("Item %d, “%s”: wrong value. Expecting “%s”; found “%s”", i, item.description, item.expected.IP, handler.IP)
		}
	}
}
