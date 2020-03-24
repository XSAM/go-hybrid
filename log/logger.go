package log

import (
	"context"

	"go.uber.org/zap"

	"github.com/XSAM/go-hybrid/consts"
)

type contextKey int

const ContextKey = contextKey(1)

const (
	// ScopeKey to distinguish scope of logs. Convenient for searching log.
	// e.g. scope: process-receipt
	//
	//   log.WithKeyValue(ctx, log.ScopeKey, "process-receipt")
	ScopeKey = "scope"
)

// Logger gets a contextual logger from current context.
// contextual logger will output common fields from context.
func Logger(ctx context.Context) *Core {
	if ctx == nil {
		return bgLogger.clone()
	}

	if ctxLogger, ok := ctx.Value(ContextKey).(*Core); ok {
		return ctxLogger
	}
	return bgLogger.clone()
}

// WithLogger add logger to context
func WithLogger(ctx context.Context, logger *Core) context.Context {
	return context.WithValue(ctx, ContextKey, logger)
}

// BgLogger return background logger
func BgLogger() *Core {
	return bgLogger
}

// SetBgLogger set background logger
func SetBgLogger(logger *Core) {
	bgLogger = logger.clone()
}

// WithTraceID attach trace id to logger
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return WithKeyValue(ctx, consts.TraceID, traceID)
}

// WithKeyValue attach key/value to logger
func WithKeyValue(ctx context.Context, key, value string) context.Context {
	var logger *Core
	if ctxLogger, ok := ctx.Value(ContextKey).(*Core); ok {
		logger = ctxLogger
	} else {
		logger = bgLogger.clone()
	}
	logger.Logger = logger.With(zap.String(key, value))

	return WithLogger(ctx, logger)
}

// WithZapOptions clones the context's Logger, applies the supplied Options.
func WithZapOptions(ctx context.Context, option ...zap.Option) context.Context {
	var logger *Core
	if ctxLogger, ok := ctx.Value(ContextKey).(*Core); ok {
		logger = ctxLogger
	} else {
		logger = bgLogger.clone()
	}
	logger.Logger = logger.WithOptions(option...)

	return WithLogger(ctx, logger)
}
