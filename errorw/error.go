package errorw

import (
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
var _ interface{ GRPCStatus() *status.Status } = (*Error)(nil)

func (e *Error) Error() string {
	return Render(e)
}

// Cause implement errors.Cause interface.
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

// APIErrorCause return the internal error's gRPC status or the root cause of the API error.
// Priority returns the internal error's gRPC status if it implements status.GRPCStatus.
// If no gRPC status can be use, then create a gRPC status with internal error.
// Implement gRPC status.GRPCStatus function.
func (e *Error) GRPCStatus() *status.Status {
	if e.Err != nil {
		if se, ok := e.Err.(interface {
			GRPCStatus() *status.Status
		}); ok {
			return se.GRPCStatus()
		}
	}
	
	st := e.APIErrorCause()
	if st != nil {
		return st
	}
	
	if e.Err != nil {
		return status.New(codes.Internal, e.Err.Error())
	}
	return nil
}

// WithAPIError append API error to error
func (e *Error) WithAPIError(apiError *status.Status) *Error {
	e.APIErrors = append(e.APIErrors, apiError)
	return e
}

// WithField append key/value to error
func (e *Error) WithField(key string, value interface{}) *Error {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}

	e.Fields[key] = value
	return e
}

// WithFields append fields to error.
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

// WithTraceID set trace id to error
func (e *Error) WithTraceID(traceID string) *Error {
	e.TraceID = traceID
	return e
}

// New create an error
func New(err error) *Error {
	return newError(err, 4)
}

func newError(err error, skip int) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		Err:   err,
		Stack: callers(skip),
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
	return New(err).WithWrap(message)
}

// NewMessage create an error with message.
func NewMessage(message string) *Error {
	return newError(errors.New(message), 4)
}

// NewMessagef create an error with message.
// It will formats according to a format specifier and returns the resulting string.
func NewMessagef(format string, args ...interface{}) *Error {
	return newError(fmt.Errorf(format, args...), 4)
}

// NewAPIError create an error and append API error.
func NewAPIError(apiError *status.Status) *Error {
	return newError(errors.New(apiError.Message()), 4).
		WithAPIError(apiError)
}
