package log

import (
	"fmt"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggingFormat(t *testing.T) {
	exitCh := make(chan int, 1)
	monkey.Patch(os.Exit, func(code int) {
		exitCh <- code
	})
	defer monkey.Unpatch(os.Exit)

	// New logger
	logger, logs := NewObservedLoggerWithLevel(zapcore.DebugLevel, zap.AddCaller())

	testCases := []struct {
		name             string
		f                func(format string, a ...interface{})
		f2               func(msg string, fields ...zap.Field)
		expectedLogLevel zapcore.Level
		expectedPanic    bool
		expectedFatal    bool
	}{
		{name: "debug", f2: logger.Debug, expectedLogLevel: zapcore.DebugLevel},
		{name: "info", f2: logger.Info, expectedLogLevel: zapcore.InfoLevel},
		{name: "warn", f2: logger.Warn, expectedLogLevel: zapcore.WarnLevel},
		{name: "error", f2: logger.Error, expectedLogLevel: zapcore.ErrorLevel},
		{name: "dpanic", f2: logger.DPanic, expectedLogLevel: zapcore.DPanicLevel},
		{name: "panic", f2: logger.Panic, expectedLogLevel: zapcore.PanicLevel, expectedPanic: true},
		{name: "fatal", f2: logger.Fatal, expectedLogLevel: zapcore.FatalLevel, expectedFatal: true},
		{name: "debugf", f: logger.Debugf, expectedLogLevel: zapcore.DebugLevel},
		{name: "infof", f: logger.Infof, expectedLogLevel: zapcore.InfoLevel},
		{name: "warnf", f: logger.Warnf, expectedLogLevel: zapcore.WarnLevel},
		{name: "errorf", f: logger.Errorf, expectedLogLevel: zapcore.ErrorLevel},
		{name: "dpanicf", f: logger.DPanicf, expectedLogLevel: zapcore.DPanicLevel},
		{name: "panicf", f: logger.Panicf, expectedLogLevel: zapcore.PanicLevel, expectedPanic: true},
		{name: "fatalf", f: logger.Fatalf, expectedLogLevel: zapcore.FatalLevel, expectedFatal: true},
	}

	format := "%s: %d"
	args := []interface{}{"foo", 42}

	message := "foo"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedPanic {
				assert.Panics(t, func() {
					if tc.f != nil {
						tc.f(format, args...)
					} else {
						tc.f2(message)
					}
				})
			} else {
				if tc.f != nil {
					tc.f(format, args...)
				} else {
					tc.f2(message)
				}
			}

			if tc.expectedFatal {
				assert.Equal(t, 1, <-exitCh)
			}

			entryList := logs.TakeAll()
			require.Equal(t, 1, len(entryList))
			l := entryList[0]
			if tc.f != nil {
				assert.Equal(t, fmt.Sprintf(format, args...), l.Message)
			} else {
				assert.Equal(t, message, l.Message)
			}
			assert.Equal(t, tc.expectedLogLevel, l.Level)
			assert.True(t, l.Caller.Defined)
			assert.NotContains(t, l.Caller.String(), "log/format.go")
		})
	}
}
