package log

import (
	"fmt"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLoggingFormat(t *testing.T) {
	exitCh := make(chan int, 1)
	monkey.Patch(os.Exit, func(code int) {
		exitCh <- code
	})
	defer monkey.Unpatch(os.Exit)

	// New logger
	logger, logs := NewObservedLoggerWithLevel(zapcore.DebugLevel)

	testCases := []struct {
		name             string
		f                func(format string, a ...interface{})
		expectedLogLevel zapcore.Level
		expectedPanic    bool
		expectedFatal    bool
	}{
		{name: "debug", f: logger.Debugf, expectedLogLevel: zapcore.DebugLevel},
		{name: "info", f: logger.Infof, expectedLogLevel: zapcore.InfoLevel},
		{name: "warn", f: logger.Warnf, expectedLogLevel: zapcore.WarnLevel},
		{name: "error", f: logger.Errorf, expectedLogLevel: zapcore.ErrorLevel},
		{name: "dpanic", f: logger.DPanicf, expectedLogLevel: zapcore.DPanicLevel},
		{name: "panic", f: logger.Panicf, expectedLogLevel: zapcore.PanicLevel, expectedPanic: true},
		{name: "fatal", f: logger.Fatalf, expectedLogLevel: zapcore.FatalLevel, expectedFatal: true},
	}

	format := "%s: %d"
	args := []interface{}{"foo", 42}

	for _, tc := range testCases {
		if tc.expectedPanic {
			assert.Panics(t, func() {
				tc.f(format, args...)
			}, tc.name)
		} else {
			tc.f(format, args...)
		}

		if tc.expectedFatal {
			assert.Equal(t, 1, <-exitCh, tc.name)
		}

		entryList := logs.TakeAll()
		require.Equal(t, 1, len(entryList), tc.name)
		l := entryList[0]
		assert.Equal(t, fmt.Sprintf(format, args...), l.Message, tc.name)
		assert.Equal(t, tc.expectedLogLevel, l.Level, tc.name)
	}
}
