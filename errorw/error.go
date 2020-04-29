package errorw

import (
	"context"
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/status"

	"github.com/XSAM/go-hybrid/trace"
)

// Error wrap error with fields and stack
type Error struct {
	Err     error
	Stack   *stack
	Wrapper []string
	Fields  map[string]interface{}

	TraceID   string
	APIErrors []*status.Status
}

type causer interface {
	Cause() error
}

// Verify interface compliance at compile time
var _ error = (*Error)(nil)
var _ stackTracer = (*Error)(nil)
var _ causer = (*Error)(nil)
var _ zapcore.ObjectMarshaler = (*Error)(nil)

func (e *Error) Error() string {
	return Render(e)
}

// Cause implement errors.Cause interface
func (e *Error) Cause() error {
	return pkgerrors.Cause(e.Err)
}

// APIErrorCause return the root cause of the API error.
func (e *Error) APIErrorCause() *status.Status {
	if len(e.APIErrors) > 0 {
		return e.APIErrors[0]
	}
	return nil
}

// WithAPIError append API error to error
func (e *Error) WithAPIError(apiError *status.Status) *Error {
	e.APIErrors = append(e.APIErrors, apiError)
	return e
}

// WithAPIError append key/value to error
func (e *Error) WithField(key string, value interface{}) *Error {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}

	e.Fields[key] = value
	return e
}

// WithAPIError append fields to error.
// Parameter fields will cover value which key is already exists.
func (e *Error) WithFields(fields map[string]interface{}) *Error {
	if e.Fields == nil {
		e.Fields = fields
	} else {
		for k, v := range fields {
			e.Fields[k] = v
		}
	}
	return e
}

// WithWrap wrap message to error
func (e *Error) WithWrap(message string) *Error {
	e.Wrapper = append(e.Wrapper, message)
	return e
}

// New a error
func New(ctx context.Context, err error) *Error {
	return newError(ctx, err, 4)
}

func newError(ctx context.Context, err error, skip int) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		TraceID: trace.GetTraceIDFromContext(ctx),
		Err:     err,
		Stack:   callers(skip),
	}
}

// Wrap wrap message
func Wrap(err error, message string) *Error {
	if err == nil {
		return nil
	}

	if val, ok := err.(*Error); ok {
		return val.WithWrap(message)
	}
	return New(context.Background(), err).WithWrap(message)
}

func NewMessage(ctx context.Context, message string) *Error {
	return newError(ctx, errors.New(message), 4)
}

func NewMessagef(ctx context.Context, format string, args ...interface{}) *Error {
	return newError(ctx, fmt.Errorf(format, args...), 4)
}
