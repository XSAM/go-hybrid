package cmdutil

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/XSAM/go-hybrid/log"
)

func TestVersionCmd(t *testing.T) {
	// Mock global logger
	logger, logs := newObservedLogger()
	log.SetBgLogger(logger)

	// Init cobra command
	cmd := cobra.Command{}
	cmd.AddCommand(VersionCmd())

	// Check Availability
	assert.Equal(t, true, cmd.HasAvailableSubCommands())

	// Run command
	_, _, err := executeCommandC(&cmd, "version")
	assert.NoError(t, err)
	assert.Contains(t, logs.All()[0].Message, "application info")
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func newObservedLogger() (*log.Core, *observer.ObservedLogs) {
	ob, logs := observer.New(zapcore.InfoLevel)
	logger := log.Core{Logger: zap.New(ob)}
	return &logger, logs
}
