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
		buf.WriteString(fmt.Sprintf(", traceID: %s", e.TraceID))
	}
	return buf.String()
}

func init() {
	Render = PlainRender
}
