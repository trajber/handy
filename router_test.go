package mux

import (
	"testing"
)

func TestAppendPath(t *testing.T) {
	p := NewPath()

	s := new(DefaultService)

	ok := p.AppendRoute("/test", s)

	t.Log(ok)
}
