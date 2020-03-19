package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogger(t *testing.T) {
	newLogger := BuildLogger(Config{
		ZapConfig: zap.NewDevelopmentConfig(),
		ZapLevel:  zapcore.InfoLevel,
	})
	ctxWithLogger := WithLogger(context.Background(), newLogger)

	testCases := []struct {
		name           string
		ctx            context.Context
		expectedLogger *Core
	}{
		{
			name:           "context is nil",
			expectedLogger: BgLogger(),
		},
		{
			name:           "context has no logger",
			ctx:            context.Background(),
			expectedLogger: BgLogger(),
		},
		{
			name:           "context has logger",
			ctx:            ctxWithLogger,
			expectedLogger: newLogger,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := Logger(tc.ctx)

			assert.Equal(t, tc.expectedLogger.Logger, logger.Logger)
		})
	}
}

func TestBgLogger(t *testing.T) {
	assert.NotNil(t, BgLogger())
}

func TestSetBgLogger(t *testing.T) {
	newLogger := BuildLogger(Config{
		ZapConfig: zap.NewDevelopmentConfig(),
		ZapLevel:  zapcore.InfoLevel,
	})
	assert.NotEqual(t, BgLogger(), newLogger)

	SetBgLogger(newLogger)
	assert.Equal(t, BgLogger(), newLogger)
}

func TestLoggerOutput(t *testing.T) {
	ctx, logs := NewContextWithObservedLogger()

	Logger(ctx).Info("testing output", zap.String("key", "value"))

	log := logs.All()[0]
	assert.Equal(t, "testing output", log.Message)
	contextMap := log.ContextMap()
	assert.Equal(t, "value", contextMap["key"])
}

func TestWithTraceID(t *testing.T) {
	ctx, logs := NewContextWithObservedLogger()
	ctx = WithTraceID(ctx, "trace_id")

	Logger(ctx).Info("testing output")

	log := logs.All()[0]
	assert.Equal(t, "testing output", log.Message)
	contextMap := log.ContextMap()
	assert.Equal(t, "trace_id", contextMap["trace_id"])
}

func TestWithKeyValue(t *testing.T) {
	ctx, logs := NewContextWithObservedLogger()

	// Set background logger
	bgLogger, bgLogs := NewObservedLogger()
	SetBgLogger(bgLogger)

	testCases := []struct {
		name string
		ctx  context.Context
		logs *observer.ObservedLogs
	}{
		{
			name: "context is nil",
		},
		{
			name: "context has no logger",
			ctx:  context.Background(),
			logs: bgLogs,
		},
		{
			name: "context has logger",
			ctx:  ctx,
			logs: logs,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ctx context.Context
			if tc.ctx == nil {
				assert.Panics(t, func() {
					ctx = WithKeyValue(tc.ctx, "foo", "bar")
				})
				return
			} else {
				ctx = WithKeyValue(tc.ctx, "foo", "bar")
			}
			Logger(ctx).Info("testing")

			log := tc.logs.All()[0]
			assert.Equal(t, "testing", log.Message)
			contextMap := log.ContextMap()
			assert.Equal(t, "bar", contextMap["foo"])
		})
	}
}

func TestWithZapOptions(t *testing.T) {
	ctx, logs := NewContextWithObservedLogger()

	// Set background logger
	bgLogger, bgLogs := NewObservedLogger()
	SetBgLogger(bgLogger)

	testCases := []struct {
		name string
		ctx  context.Context
		logs *observer.ObservedLogs
	}{
		{
			name: "context is nil",
		},
		{
			name: "context has no logger",
			ctx:  context.Background(),
			logs: bgLogs,
		},
		{
			name: "context has logger",
			ctx:  ctx,
			logs: logs,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ctx context.Context
			if tc.ctx == nil {
				assert.Panics(t, func() {
					ctx = WithZapOptions(tc.ctx, zap.Fields(zap.String("foo", "bar")))
				})
				return
			} else {
				ctx = WithZapOptions(tc.ctx, zap.Fields(zap.String("foo", "bar")))
			}
			Logger(ctx).Info("testing")

			log := tc.logs.All()[0]
			assert.Equal(t, "testing", log.Message)
			contextMap := log.ContextMap()
			assert.Equal(t, "bar", contextMap["foo"])
		})
	}
}

func NewObservedLogger() (*Core, *observer.ObservedLogs) {
	ob, logs := observer.New(zapcore.InfoLevel)
	logger := Core{Logger: zap.New(ob)}
	return &logger, logs
}

func NewContextWithObservedLogger() (context.Context, *observer.ObservedLogs) {
	logger, logs := NewObservedLogger()
	ctxWithLogger := WithLogger(context.Background(), logger)
	return ctxWithLogger, logs
}
