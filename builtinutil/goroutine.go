package builtinutil

import (
	"context"

	"go.uber.org/zap"
)

func WrappedGo(ctx context.Context, callback func(), options ...zap.Option) {
	defer RecoveryWithContext(ctx, nil, options...)
	callback()
}
