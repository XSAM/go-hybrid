package zapfield

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/XSAM/go-hybrid/environment"
	"github.com/XSAM/go-hybrid/errorw"
	"github.com/XSAM/go-hybrid/log"
)

func TestStack(t *testing.T) {
	logger, logs := newObservedLogger()

	logger.Info("testing", Stack())

	l := logs.All()[0]
	assert.Equal(t, "testing", l.Message)
	contextMap := l.ContextMap()
	assert.Contains(t, contextMap["stack"], "go-hybrid/log/zapfield/field_test.go")
}

func TestError(t *testing.T) {
	logger, logs := newObservedLogger()

	normalError := errors.New("error")
	errorwError := errorw.New(normalError)

	testCases := []struct {
		mode        environment.ModeType
		interaction environment.LogStyleType
		err         error
	}{
		// Normal error
		{
			mode:        environment.ModeDevelopment,
			interaction: environment.LogStyleJSON,
			err:         normalError,
		},
		{
			mode:        environment.ModeDevelopment,
			interaction: environment.LogStyleText,
			err:         normalError,
		},
		{
			mode:        environment.ModeProduction,
			interaction: environment.LogStyleJSON,
			err:         normalError,
		},
		{
			mode:        environment.ModeProduction,
			interaction: environment.LogStyleText,
			err:         normalError,
		},

		// errorw error
		{
			mode:        environment.ModeDevelopment,
			interaction: environment.LogStyleJSON,
			err:         errorwError,
		},
		{
			mode:        environment.ModeDevelopment,
			interaction: environment.LogStyleText,
			err:         errorwError,
		},
		{
			mode:        environment.ModeProduction,
			interaction: environment.LogStyleJSON,
			err:         errorwError,
		},
		{
			mode:        environment.ModeProduction,
			interaction: environment.LogStyleText,
			err:         errorwError,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s-%s-%T", tc.mode, tc.interaction, tc.err), func(t *testing.T) {
			environment.LogStyle = tc.interaction
			environment.Mode = tc.mode
			logger.Info("testing", Error(tc.err))

			l := logs.All()[0]
			assert.Equal(t, "testing", l.Message)
			contextMap := l.ContextMap()
			assert.NotEmpty(t, contextMap["error"])
		})
	}
}

func TestNilError(t *testing.T) {
	logger, logs := newObservedLogger()
	logger.Info("testing", Error(nil))
	
	l := logs.All()[0]
	assert.Equal(t, "testing", l.Message)
	contextMap := l.ContextMap()
	assert.Empty(t, contextMap["error"])
}

func newObservedLogger() (*log.Core, *observer.ObservedLogs) {
	ob, logs := observer.New(zapcore.InfoLevel)
	logger := log.Core{Logger: zap.New(ob)}
	return &logger, logs
}
