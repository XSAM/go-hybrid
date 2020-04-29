package errorw

import (
	"bytes"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MarshalLogObject is an implementation of `zapcore.ObjectMarshaler` interface
func (e *Error) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	// Trace ID
	if e.TraceID != "" {
		enc.AddString("trace_id", e.TraceID)
	}

	// Error message
	var buf bytes.Buffer
	for i := len(e.Wrapper) - 1; i >= 0; i-- {
		buf.WriteString(fmt.Sprintf("%s: ", e.Wrapper[i]))
	}
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		buf.WriteString("nil")
	}
	enc.AddString("msg", buf.String())

	// Stack
	enc.AddString("stack", fmt.Sprintf("%+v", e.Stack))

	// Fields
	if len(e.Fields) > 0 {
		field := zap.Any("fields", e.Fields)
		field.AddTo(enc)
	}
	return nil
}
