package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPresetEnvironments(t *testing.T) {
	testCases := []struct {
		name                    string
		function                func()
		expectedModeType        ModeType
		expectedInteractionType InteractionType
	}{
		{
			name:                    "development mode",
			function:                DevelopmentMode,
			expectedModeType:        DEVELOPMENT_MODE,
			expectedInteractionType: NORMAL_INTERACTIION,
		},
		{
			name:                    "production mode",
			function:                ProductionMode,
			expectedModeType:        PRODUCTION_MODE,
			expectedInteractionType: NORMAL_INTERACTIION,
		},
		{
			name:                    "CLI tool development mode",
			function:                CLIToolDevelopmentMode,
			expectedModeType:        DEVELOPMENT_MODE,
			expectedInteractionType: CLI_INTERACTION,
		},
		{
			name:                    "CLI tool production mode",
			function:                CLIToolProductionMode,
			expectedModeType:        PRODUCTION_MODE,
			expectedInteractionType: CLI_INTERACTION,
		},
	}

	for _, tc := range testCases {
		tc.function()
		assert.Equal(t, tc.expectedModeType, Mode)
		assert.Equal(t, tc.expectedInteractionType, Interaction)
	}
}
