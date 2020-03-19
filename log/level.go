package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Levels map[string]*zap.AtomicLevel

// NewLogLevels return a new log levels
func NewLogLevels(defaultLevel zapcore.Level) *Levels {
	d := zap.NewAtomicLevelAt(defaultLevel)
	level := make(Levels)
	level["default"] = &d

	return &level
}

// GetLevels return a global log levels
func GetLevels() *Levels {
	return logLevels
}

func (l *Levels) Get() *zap.AtomicLevel {
	return l.GetWithScope("default")
}

func (l *Levels) Set(level zapcore.Level) {
	l.SetWithScope("default", level)
}

func (l *Levels) GetWithScope(scope string) *zap.AtomicLevel {
	return (*l)[scope]
}

func (l *Levels) SetWithScope(scope string, level zapcore.Level) {
	if _, ok := (*l)[scope]; ok {
		(*l)[scope].SetLevel(level)
	} else {
		al := zap.NewAtomicLevelAt(level)
		(*l)[scope] = &al
	}
}
