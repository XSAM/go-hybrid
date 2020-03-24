package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPresetConfig(t *testing.T) {
	presetList := []func() Config{
		DevelopmentAndJSONConfig,
		ProductionAndJSONConfig,
		DevelopmentAndTextConfig,
		ProductionAndTextConfig,
	}

	for _, preset := range presetList {
		logger := BuildLogger(preset())
		assert.NotNil(t, logger)
	}
}

func TestBuildLogger(t *testing.T) {
	testCases := []struct {
		name          string
		config        Config
		expectedPanic bool
	}{
		{
			name:   "valid config",
			config: DevelopmentAndJSONConfig(),
		},
		{
			name:          "invalid config",
			config:        Config{},
			expectedPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedPanic {
				assert.Panics(t, func() {
					BuildLogger(tc.config)
				})
			} else {
				logger := BuildLogger(tc.config)
				assert.NotNil(t, logger)
			}
		})
	}
}
