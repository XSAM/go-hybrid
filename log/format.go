package log

import "fmt"

// Debugf formats according to a format specifier and logs the formatted message at DebugLevel.
func (c *Core) Debugf(format string, v ...interface{}) {
	c.Debug(fmt.Sprintf(format, v...))
}

// Infof formats according to a format specifier and logs the formatted message at InfoLevel.
func (c *Core) Infof(format string, v ...interface{}) {
	c.Info(fmt.Sprintf(format, v...))
}

// Warnf formats according to a format specifier and logs the formatted message at WarnLevel.
func (c *Core) Warnf(format string, v ...interface{}) {
	c.Warn(fmt.Sprintf(format, v...))
}

// Errorf formats according to a format specifier and logs the formatted message at ErrorLevel.
func (c *Core) Errorf(format string, v ...interface{}) {
	c.Error(fmt.Sprintf(format, v...))
}

// DPanicf formats according to a format specifier and logs the formatted message at DPanicLevel.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (c *Core) DPanicf(format string, v ...interface{}) {
	c.DPanic(fmt.Sprintf(format, v...))
}

// Panicf formats according to a format specifier and logs the formatted message at PanicLevel.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (c *Core) Panicf(format string, v ...interface{}) {
	c.Panic(fmt.Sprintf(format, v...))
}

// Fatalf formats according to a format specifier and logs the formatted message at FatalLevel.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (c *Core) Fatalf(format string, v ...interface{}) {
	c.Fatal(fmt.Sprintf(format, v...))
}
