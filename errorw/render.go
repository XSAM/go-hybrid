package errorw

import (
	"bytes"
	"fmt"
)

var Render func(e *Error) string

func PlainRender(e *Error) string {
	var buf bytes.Buffer

	for i := len(e.Wrapper) - 1; i >= 0; i-- {
		buf.WriteString(fmt.Sprintf("%s: ", e.Wrapper[i]))
	}
	buf.WriteString(e.Err.Error())

	if e.TraceID != "" {
		buf.WriteString(fmt.Sprintf(". traceID: %s", e.TraceID))
	}

	if len(e.Fields) > 0 {
		buf.WriteString(fmt.Sprintf(". fields:"))
		for k, v := range e.Fields {
			buf.WriteString(fmt.Sprintf(" %s:%+v", k, v))
		}
	}
	return buf.String()
}

func init() {
	Render = PlainRender
}
