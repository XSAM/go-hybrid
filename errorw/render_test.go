package errorw

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

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

	assert.Equal(t, "test2: test: foo. traceID: trace_id. fields: foo:bar struct_field:{Foo:foo Bar:[bar1 bar2]}", PlainRender(err))
}
