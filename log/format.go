package log

import (
	"fmt"

	"go.uber.org/zap"
)

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (c *Core) Debug(msg string, fields ...zap.Field) {
	c.Logger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (c *Core) Info(msg string, fields ...zap.Field) {
	c.Logger.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (c *Core) Warn(msg string, fields ...zap.Field) {
	c.Logger.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (c *Core) Error(msg string, fields ...zap.Field) {
	c.Logger.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (c *Core) DPanic(msg string, fields ...zap.Field) {
	c.Logger.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (c *Core) Panic(msg string, fields ...zap.Field) {
	c.Logger.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (c *Core) Fatal(msg string, fields ...zap.Field) {
	c.Logger.Fatal(msg, fields...)
}

// Debugf formats according to a format specifier and logs the formatted message at DebugLevel.
func (c *Core) Debugf(format string, v ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(format, v...))
}

// Infof formats according to a format specifier and logs the formatted message at InfoLevel.
func (c *Core) Infof(format string, v ...interface{}) {
	c.Logger.Info(fmt.Sprintf(format, v...))
}

// Warnf formats according to a format specifier and logs the formatted message at WarnLevel.
func (c *Core) Warnf(format string, v ...interface{}) {
	c.Logger.Warn(fmt.Sprintf(format, v...))
}

// Errorf formats according to a format specifier and logs the formatted message at ErrorLevel.
func (c *Core) Errorf(format string, v ...interface{}) {
	c.Logger.Error(fmt.Sprintf(format, v...))
}

// DPanicf formats according to a format specifier and logs the formatted message at DPanicLevel.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (c *Core) DPanicf(format string, v ...interface{}) {
	c.Logger.DPanic(fmt.Sprintf(format, v...))
}

// Panicf formats according to a format specifier and logs the formatted message at PanicLevel.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (c *Core) Panicf(format string, v ...interface{}) {
	c.Logger.Panic(fmt.Sprintf(format, v...))
}

// Fatalf formats according to a format specifier and logs the formatted message at FatalLevel.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (c *Core) Fatalf(format string, v ...interface{}) {
	c.Logger.Fatal(fmt.Sprintf(format, v...))
}
