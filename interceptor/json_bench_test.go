package interceptor

import (
	"br/tests"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"br/Godeps/_workspace/src/github.com/gustavo-hms/handy"
)

type mockJSONHandler struct {
	IntrospectorCompliant

	req     *http.Request
	resp    http.ResponseWriter
	Request struct {
		Um   int
		Dois string
		Três struct {
			Quatro []int
		}
	} `request:"put"`
	Response struct {
		Cinco int
		Seis  string
		Sete  []bool
	} `response:"get"`
}

func (m mockJSONHandler) Req() *http.Request {
	return m.req
}

func (m mockJSONHandler) ResponseWriter() http.ResponseWriter {
	return m.resp
}

func TestJSONBefore(t *testing.T) {
	json := `
	{
		"um": 1,
		"dois": "dois",
		"três": {
			"quatro": [1, 2, 3, 4, 5]
		}
	}
	`

	req, err := http.NewRequest("PUT", "/", strings.NewReader(json))

	if err != nil {
		t.Fatal(err)
	}

	handler := &mockJSONHandler{req: req}
	i := NewIntrospector(handler)
	i.Before()
	u := NewJSONCodec(handler)
	status := u.Before()

	if status != 0 {
		t.Errorf("Wrong status code. Expecting “0”; found “%d”", status)
	}

	expected := struct {
		Um   int
		Dois string
		Três struct {
			Quatro []int
		}
	}{
		Um:   1,
		Dois: "dois",
		Três: struct {
			Quatro []int
		}{
			Quatro: []int{1, 2, 3, 4, 5},
		},
	}

	if !reflect.DeepEqual(expected, handler.Request) {
		t.Error("Wrong request")
		t.Log(tests.Diff(expected, handler.Request))
	}
}

func TestJSONAfter(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := &mockJSONHandler{
		req:  r,
		resp: w,
		Response: struct {
			Cinco int
			Seis  string
			Sete  []bool
		}{
			Cinco: 5,
			Seis:  "seis",
			Sete:  []bool{true, true, false},
		},
	}

	i := NewIntrospector(handler)
	i.Before()
	u := NewJSONCodec(handler)
	status := u.After(http.StatusOK)

	if status != http.StatusOK {
		t.Errorf("Wrong status code. Expecting “200”; found “%d”", status)
	}

	expected := `{"Cinco":5,"Seis":"seis","Sete":[true,true,false]}` + "\n"

	if w.Body.String() != expected {
		t.Errorf("Wrong response. Expecting “%s”; found “%s”", expected, w.Body.String())
	}
}

var (
	payload = `
{
	"name":"foo",
	"id":10
}`
)

type TestStruct struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func BenchmarkDecodeJSON(b *testing.B) {
	req, err := http.NewRequest("GET", "/", strings.NewReader(payload))
	if err != nil {
		b.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := new(struct {
		handy.DefaultHandler
		IntrospectorCompliant
		Request TestStruct `request:"get"`
	})
	handy.SetHandlerInfo(handler, w, req, nil)
	codec := NewJSONCodec(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		codec.Before()
		if handler.Request.ID != 10 {
			b.Fail()
		}
	}
}
