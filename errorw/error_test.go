package errorw

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/XSAM/go-hybrid/trace"
)

func TestNew(t *testing.T) {
	tcErr := errors.New("testing")

	testCases := []struct {
		name string
		ctx  context.Context
		err  error
	}{
		{
			name: "context is nil, error is nil",
		},
		{
			name: "context is not nil, error is nil",
			ctx:  context.Background(),
		},
		{
			name: "context and error not nil",
			ctx:  context.Background(),
		},
		{
			name: "context with trace id",
			ctx:  trace.SetTraceIDToContext(context.Background(), "trace_id"),
			err:  tcErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := New(tc.ctx, tc.err)

			if tc.err == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, trace.GetTraceIDFromContext(tc.ctx), err.TraceID)
				assert.NotNil(t, err.Stack)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "error is nil",
		},
		{
			name: "error type is errorw.Error",
			err:  New(context.Background(), errors.New("testing")),
		},
		{
			name: "error type is not errorw.Error",
			err:  errors.New("testing"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Wrap(tc.err, "foo")
			if tc.err == nil {
				assert.Nil(t, err)
			} else {
				assert.Len(t, err.Wrapper, 1)
				assert.Equal(t, "foo", err.Wrapper[0])
			}
		})
	}
}

func TestWrap2(t *testing.T) {
	err := errors.New("testing")

	messages := []string{"foo", "bar", "ver"}
	for _, v := range messages {
		err = Wrap(err, v)
	}

	e := err.(*Error)
	assert.Equal(t, messages, e.Wrapper)
}

func TestError_Error(t *testing.T) {
	err := New(trace.SetTraceIDToContext(context.Background(), "trace_id"), errors.New("foo"))
	err = err.WithField("foo", "bar").
		WithWrap("test").
		WithWrap("test2")

	assert.Equal(t, "test2: test: foo. traceID: trace_id. fields: foo:bar", err.Error())
}

func TestError_Cause(t *testing.T) {
	err := errors.New("foo")
	e := New(context.Background(), err)

	assert.Equal(t, err, e.Cause())
}

func TestError_APIErrorCause(t *testing.T) {
	// Have value
	err := &Error{APIErrors: []*status.Status{
		status.New(codes.Internal, "foo"),
		status.New(codes.Internal, "bar"),
	}}
	assert.Equal(t, "foo", err.APIErrorCause().Message())

	// Empty
	err = &Error{}
	assert.Nil(t, err.APIErrorCause())
}

func TestError_WithAPIError(t *testing.T) {
	err := New(context.Background(), errors.New("foo")).
		WithAPIError(status.New(codes.Internal, "foo")).
		WithAPIError(status.New(codes.Internal, "bar"))

	assert.Equal(t, "foo", err.APIErrorCause().Message())
}

func TestError_WithField(t *testing.T) {
	err := New(context.Background(), errors.New("foo")).
		WithField("foo", "foo").
		WithField("bar", "bar")

	assert.Equal(t, map[string]interface{}{
		"foo": "foo",
		"bar": "bar",
	}, err.Fields)
}

func TestError_WithFields(t *testing.T) {
	err := New(context.Background(), errors.New("foo")).WithFields(map[string]interface{}{
		"foo": "foo",
		"bar": "bar",
	}).WithFields(map[string]interface{}{
		"foo":  "foo2",
		"bar2": "bar",
	})

	assert.Equal(t, map[string]interface{}{
		"foo":  "foo2",
		"bar":  "bar",
		"bar2": "bar",
	}, err.Fields)
}

func TestError_WithWrap(t *testing.T) {
	err := New(context.Background(), errors.New("foo")).
		WithWrap("foo").WithWrap("bar")

	assert.Equal(t, []string{"foo", "bar"}, err.Wrapper)
}
