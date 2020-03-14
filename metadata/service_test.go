package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppName(t *testing.T) {
	// Before setup app name
	assert.Equal(t, "", AppName())

	// After setup app name
	SetAppName("foo")
	assert.Equal(t, "foo", AppName())
}

func TestRuntimeID(t *testing.T) {
	// Mock runtime id
	appInfo.RuntimeID = "foo"

	assert.Equal(t, "foo", RuntimeID())
}

func TestAppInfo(t *testing.T) {
	info := AppInfo()

	// Runtime ID not empty
	assert.NotEqual(t, "", info.RuntimeID)
}
