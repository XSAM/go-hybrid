package log

import (
	"go.uber.org/zap"
)

type Core struct {
	*zap.Logger
}
