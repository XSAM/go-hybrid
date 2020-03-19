package zapfield

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
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
	errorwError := errorw.New(nil, normalError)

	testCases := []struct {
		mode        environment.ModeType
		interaction environment.InteractionType
		err         error
	}{
		// Normal error
		{
			mode:        environment.DEVELOPMENT_MODE,
			interaction: environment.NORMAL_INTERACTIION,
			err:         normalError,
		},
		{
			mode:        environment.DEVELOPMENT_MODE,
			interaction: environment.CLI_INTERACTION,
			err:         normalError,
		},
		{
			mode:        environment.PRODUCTION_MODE,
			interaction: environment.NORMAL_INTERACTIION,
			err:         normalError,
		},
		{
			mode:        environment.PRODUCTION_MODE,
			interaction: environment.CLI_INTERACTION,
			err:         normalError,
		},

		// errorw error
		{
			mode:        environment.DEVELOPMENT_MODE,
			interaction: environment.NORMAL_INTERACTIION,
			err:         errorwError,
		},
		{
			mode:        environment.DEVELOPMENT_MODE,
			interaction: environment.CLI_INTERACTION,
			err:         errorwError,
		},
		{
			mode:        environment.PRODUCTION_MODE,
			interaction: environment.NORMAL_INTERACTIION,
			err:         errorwError,
		},
		{
			mode:        environment.PRODUCTION_MODE,
			interaction: environment.CLI_INTERACTION,
			err:         errorwError,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s-%s-%T", tc.mode, tc.interaction, tc.err), func(t *testing.T) {
			environment.Interaction = tc.interaction
			environment.Mode = tc.mode
			logger.Info("testing", Error(tc.err))

			l := logs.All()[0]
			assert.Equal(t, "testing", l.Message)
			contextMap := l.ContextMap()
			assert.NotEmpty(t, contextMap["error"])
		})
	}
}

func newObservedLogger() (*log.Core, *observer.ObservedLogs) {
	ob, logs := observer.New(zapcore.InfoLevel)
	logger := log.Core{Logger: zap.New(ob)}
	return &logger, logs
}
