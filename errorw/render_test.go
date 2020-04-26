package errorw

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/XSAM/go-hybrid/trace"
)

func TestPlainRender(t *testing.T) {
	structField := struct {
		Foo string
		Bar []string
	}{
		Foo: "foo",
		Bar: []string{"bar1", "bar2"},
	}

	err := New(trace.SetTraceIDToContext(context.Background(), "trace_id"), errors.New("foo"))
	err = err.WithField("foo", "bar").
		WithField("struct_field", structField).
		WithWrap("test").
		WithWrap("test2")

	expectedResult1 := "test2: test: foo. traceID: trace_id. fields: foo:bar struct_field:{Foo:foo Bar:[bar1 bar2]}"
	expectedResult2 := "test2: test: foo. traceID: trace_id. fields: struct_field:{Foo:foo Bar:[bar1 bar2]} foo:bar"
	result := PlainRender(err)

	if result != expectedResult1 && result != expectedResult2 {
		t.Fatal("unexpected result")
	}
}
