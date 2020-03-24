package log

import (
	"go.uber.org/zap"
)

type Core struct {
	*zap.Logger
}

func (c *Core) clone() *Core {
	newCore := *c
	return &newCore
}
