package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogLevel(t *testing.T) {
	level := NewLogLevels(zapcore.WarnLevel)

	l := zap.NewAtomicLevelAt(zapcore.WarnLevel)
	assert.Equal(t, &Levels{"default": &l}, level)
}

func TestGetLevels(t *testing.T) {
	assert.NotNil(t, GetLevels())
}

type LevelsTestSuite struct {
	suite.Suite
	levels       *Levels
	defaultLevel zap.AtomicLevel
}

func (suite *LevelsTestSuite) SetupSuite() {
	suite.levels = NewLogLevels(zapcore.WarnLevel)
	suite.defaultLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
}

func (suite *LevelsTestSuite) TestGet() {
	t := suite.T()
	result := suite.levels.Get()

	assert.Equal(t, &suite.defaultLevel, result)
}

func (suite *LevelsTestSuite) TestSet() {
	t := suite.T()

	newLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	suite.levels.Set(zapcore.DebugLevel)

	assert.Equal(t, &newLevel, suite.levels.Get())
}

func (suite *LevelsTestSuite) TestGetAndSetWithScope() {
	t := suite.T()
	testCases := []struct {
		name  string
		scope string
		level zapcore.Level
	}{
		{
			name:  "setup new scope",
			scope: "a",
			level: zapcore.FatalLevel,
		},
		{
			name:  "setup exists scope",
			scope: "a",
			level: zapcore.DebugLevel,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newLevel := zap.NewAtomicLevelAt(tc.level)
			suite.levels.SetWithScope(tc.scope, tc.level)

			assert.Equal(t, &newLevel, suite.levels.GetWithScope(tc.scope))
		})
	}
	result := suite.levels.Get()

	assert.Equal(t, &suite.defaultLevel, result)
}

func TestLevelsTestSuite(t *testing.T) {
	suite.Run(t, new(LevelsTestSuite))
}
