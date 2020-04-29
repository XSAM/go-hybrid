package errorw

import (
	"errors"
	"testing"
)

func TestPlainRender(t *testing.T) {
	structField := struct {
		Foo string
		Bar []string
	}{
		Foo: "foo",
		Bar: []string{"bar1", "bar2"},
	}

	err := New(errors.New("foo"))
	err = err.WithField("foo", "bar").
		WithField("struct_field", structField).
		WithWrap("test").
		WithWrap("test2").
		WithTraceID("trace_id")

	expectedResult1 := "test2: test: foo. traceID: trace_id. fields: foo:bar struct_field:{Foo:foo Bar:[bar1 bar2]}"
	expectedResult2 := "test2: test: foo. traceID: trace_id. fields: struct_field:{Foo:foo Bar:[bar1 bar2]} foo:bar"
	result := PlainRender(err)

	if result != expectedResult1 && result != expectedResult2 {
		t.Fatal("unexpected result")
	}
}
