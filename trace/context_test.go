package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/XSAM/go-hybrid/consts"
)

func TestSetTraceIDToContext(t *testing.T) {
	ctx := SetTraceIDToContext(context.Background(), "trace_id")

	assert.Equal(t, "trace_id", ctx.Value(consts.TraceID).(string))
}

func TestGetTraceIDFromContext(t *testing.T) {
	testCases := []struct {
		name            string
		ctx             context.Context
		expectedTraceID string
	}{
		{
			name: "context is nil",
		},
		{
			name: "context has no trace id",
			ctx:  context.Background(),
		},
		{
			name:            "context has trace id",
			ctx:             SetTraceIDToContext(context.Background(), "trace_id"),
			expectedTraceID: "trace_id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetTraceIDFromContext(tc.ctx)
			assert.Equal(t, tc.expectedTraceID, result)
		})
	}
}
