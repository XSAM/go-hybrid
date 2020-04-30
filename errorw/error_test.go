package errorw

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNew(t *testing.T) {
	tcErr := errors.New("testing")

	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "error is nil",
		},
		{
			name: "err is not nil",
			err:  tcErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := New(tc.err)

			if tc.err == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
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
			err:  New(errors.New("testing")),
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
	err := New(errors.New("foo"))
	err = err.WithField("foo", "bar").
		WithWrap("test").
		WithWrap("test2").
		WithTraceID("trace_id")

	assert.Equal(t, "test2: test: foo. traceID: trace_id. fields: foo:bar", err.Error())
}

func TestError_Cause(t *testing.T) {
	err := errors.New("foo")
	e := New(err)

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

func TestError_GRPCStatus(t *testing.T) {
	// Have value
	err := &Error{APIErrors: []*status.Status{
		status.New(codes.Internal, "foo"),
		status.New(codes.Internal, "bar"),
	}}
	assert.Equal(t, "foo", err.GRPCStatus().Message())

	// Empty
	err = &Error{}
	assert.Nil(t, err.GRPCStatus())
}

func TestError_WithAPIError(t *testing.T) {
	err := New(errors.New("foo")).
		WithAPIError(status.New(codes.Internal, "foo")).
		WithAPIError(status.New(codes.Internal, "bar"))

	assert.Equal(t, "foo", err.APIErrorCause().Message())
}

func TestError_WithField(t *testing.T) {
	err := New(errors.New("foo")).
		WithField("foo", "foo").
		WithField("bar", "bar")

	assert.Equal(t, map[string]interface{}{
		"foo": "foo",
		"bar": "bar",
	}, err.Fields)
}

func TestError_WithFields(t *testing.T) {
	err := New(errors.New("foo")).WithFields(map[string]interface{}{
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
	err := New(errors.New("foo")).
		WithWrap("foo").WithWrap("bar")

	assert.Equal(t, []string{"foo", "bar"}, err.Wrapper)
}

func TestError_WithTraceID(t *testing.T) {
	err := NewMessage("foo").WithTraceID("id")

	assert.Equal(t, "foo. traceID: id", err.Error())
	assert.Equal(t, "id", err.TraceID)
}

func TestNewMessage(t *testing.T) {
	err := NewMessage("foo")

	assert.Equal(t, "foo", err.Error())
}

func TestNewMessagef(t *testing.T) {
	format := "%s: %d"
	args := []interface{}{"foo", 42}

	err := NewMessagef(format, args...)

	assert.Equal(t, fmt.Sprintf(format, args...), err.Error())
}

func TestNewAPIError(t *testing.T) {
	apiErr := status.New(codes.Internal, "test")
	err := NewAPIError(apiErr)

	assert.Error(t, err)
	assert.Equal(t, "test", err.Err.Error())
	assert.Equal(t, "test", err.Error())
	assert.Len(t, err.APIErrors, 1)
	assert.Equal(t, apiErr, err.APIErrorCause())
	assert.Equal(t, apiErr, err.GRPCStatus())
}
