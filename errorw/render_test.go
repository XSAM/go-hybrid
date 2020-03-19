package errorw

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/XSAM/go-hybrid/trace"
)

func TestPlainRender(t *testing.T) {
	err := New(trace.SetTraceIDToContext(context.Background(), "trace_id"), errors.New("foo"))
	err = err.WithField("foo", "bar").
		WithWrap("test").
		WithWrap("test2")

	assert.Equal(t, "test2: test: foo, traceID: trace_id", PlainRender(err))
}
