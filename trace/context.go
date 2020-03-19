package trace

import (
	"context"

	"github.com/XSAM/go-hybrid/consts"
)

func GetTraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if val, ok := ctx.Value(consts.TraceID).(string); ok {
		return val
	}
	return ""
}

func SetTraceIDToContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, consts.TraceID, traceID)
}
