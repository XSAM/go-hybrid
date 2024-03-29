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

func TestWrapf(t *testing.T) {
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
			err := Wrapf(tc.err, "%s", "foo")
			if tc.err == nil {
				assert.Nil(t, err)
			} else {
				assert.Len(t, err.Wrapper, 1)
				assert.Equal(t, "foo", err.Wrapper[0])
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	err := New(errors.New("foo"))
	err = err.WithField("foo", "bar").
		WithWrap("test").
		WithWrap("test2")

	assert.Equal(t, "test2: test: foo. fields: foo:bar", err.Error())
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
	testCases := []struct {
		name               string
		err                *Error
		expectedGRPCStatus *status.Status
	}{
		{
			name: "have API error",
			err: &Error{APIErrors: []*status.Status{
				status.New(codes.Internal, "foo"),
				status.New(codes.Internal, "bar"),
			}},
			expectedGRPCStatus: status.New(codes.Internal, "foo"),
		},
		{
			name:               "empty",
			err:                &Error{},
			expectedGRPCStatus: nil,
		},
		{
			name:               "internal err implement gRPC status.GRPCStatus",
			err:                &Error{Err: status.New(codes.Internal, "root").Err()},
			expectedGRPCStatus: status.New(codes.Internal, "root"),
		},
		{
			name: "internal err implement gRPC status.GRPCStatus. Error also contain API error",
			err: &Error{
				Err: status.New(codes.Internal, "root").Err(),
				APIErrors: []*status.Status{
					status.New(codes.Internal, "foo"),
					status.New(codes.Internal, "bar"),
				}},
			expectedGRPCStatus: status.New(codes.Internal, "root"),
		},
		{
			name: "if no gRPC status can be use, then create a gRPC status with internal error",
			err: &Error{
				Err: errors.New("error"),
			},
			expectedGRPCStatus: status.New(codes.Internal, "error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.err.GRPCStatus()
			assert.Equal(t, tc.expectedGRPCStatus, result)
		})
	}
}

func TestError_WithAPIError(t *testing.T) {
	err := New(errors.New("foo")).
		WithAPIError(status.New(codes.Internal, "foo")).
		WithAPIError(status.New(codes.Internal, "bar"))

	assert.Equal(t, "foo", err.APIErrorCause().Message())
	
	assert.Nil(t, (*Error)(nil).WithAPIError(status.New(codes.Internal, "foo")))
}

func TestError_WithField(t *testing.T) {
	err := New(errors.New("foo")).
		WithField("foo", "foo").
		WithField("bar", "bar")

	assert.Equal(t, map[string]interface{}{
		"foo": "foo",
		"bar": "bar",
	}, err.Fields)
	
	assert.Nil(t, (*Error)(nil).WithField("foo", "foo"))
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
	
	assert.Nil(t, (*Error)(nil).WithFields(map[string]interface{}{
		"foo": "foo",
		"bar": "bar",
	}))
}

func TestError_WithWrap(t *testing.T) {
	err := New(errors.New("foo")).
		WithWrap("foo").WithWrap("bar")

	assert.Equal(t, []string{"foo", "bar"}, err.Wrapper)
	
	assert.Nil(t, (*Error)(nil).WithWrap("foo"))
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
