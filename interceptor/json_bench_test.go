package interceptor_test

import (
	"handy"
	"handy/interceptor"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

type jsonHandler struct {
	handy.ProtoHandler
	interceptor.JSONCodecAPI

	RequestData struct {
		Um   int
		Dois string
		Três struct {
			Quatro []int
		}
	} `request:"put"`
	ResponseData struct {
		Cinco int
		Seis  string
		Sete  []bool
	} `response:"get"`
}

func newJSONHandler(ctx handy.Context, handler *jsonHandler) (handy.Handler, handy.Interceptor) {
	intro := interceptor.NewIntrospector(nil, handler)
	json := interceptor.NewJSONCodec(intro)
	handler.JSONCodecAPI = json

	return handler, json
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

	request, err := http.NewRequest("PUT", "/", strings.NewReader(json))

	if err != nil {
		t.Fatal(err)
	}

	ctx := handy.Context{Request: request}
	handler := &jsonHandler{}
	_, intercept := newJSONHandler(ctx, handler)
	// The framework sets the context automatically, but here in the
	// tests we need to do it manually
	handler.SetContext(ctx)
	intercept.SetContext(ctx)
	status := intercept.Before()

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

	if !reflect.DeepEqual(expected, handler.RequestData) {
		t.Error("Wrong request")
		t.Log(pretty.Compare(expected, handler.RequestData))
	}
}

func TestJSONAfter(t *testing.T) {
	handler := &jsonHandler{
		ResponseData: struct {
			Cinco int
			Seis  string
			Sete  []bool
		}{
			Cinco: 5,
			Seis:  "seis",
			Sete:  []bool{true, true, false},
		},
	}

	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	writer := httptest.NewRecorder()
	ctx := handy.Context{Request: request, ResponseWriter: writer}
	_, intercept := newJSONHandler(ctx, handler)
	handler.SetContext(ctx)
	intercept.SetContext(ctx)
	status := intercept.After(http.StatusOK)

	if status != http.StatusOK {
		t.Errorf("Wrong status code. Expecting “200”; found “%d”", status)
	}

	expected := `{"Cinco":5,"Seis":"seis","Sete":[true,true,false]}`

	if writer.Body.String() != expected {
		t.Errorf("Wrong response. Expecting “%s”; found “%s”", expected, writer.Body.String())
	}
}
