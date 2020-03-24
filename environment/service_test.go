package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPresetEnvironments(t *testing.T) {
	testCases := []struct {
		name                 string
		functions            []func()
		expectedModeType     ModeType
		expectedLogStyleType LogStyleType
	}{
		{
			name:                 "development mode with json log style",
			functions:            []func(){DevelopmentMode, JSONLogStyle},
			expectedModeType:     ModeDevelopment,
			expectedLogStyleType: LogStyleJSON,
		},
		{
			name:                 "production mode with json log style",
			functions:            []func(){ProductionMode, JSONLogStyle},
			expectedModeType:     ModeProduction,
			expectedLogStyleType: LogStyleJSON,
		},
		{
			name:                 "development mode with text log style",
			functions:            []func(){DevelopmentMode, TextLogStyle},
			expectedModeType:     ModeDevelopment,
			expectedLogStyleType: LogStyleText,
		},
		{
			name:                 "production mode with text log style",
			functions:            []func(){ProductionMode, TextLogStyle},
			expectedModeType:     ModeProduction,
			expectedLogStyleType: LogStyleText,
		},
	}

	for _, tc := range testCases {
		for _, v := range tc.functions {
			v()
		}
		assert.Equal(t, tc.expectedModeType, Mode)
		assert.Equal(t, tc.expectedLogStyleType, LogStyle)
	}
}
