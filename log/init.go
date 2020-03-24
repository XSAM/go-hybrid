package log

import "go.uber.org/zap"

var logLevels *Levels
var bgLogger *Core

func init() {
	logLevels = NewLogLevels(zap.InfoLevel)

	config := ProductionAndJSONConfig()
	BuildAndSetBgLogger(config)
}

func BuildAndSetBgLogger(config Config) {
	SetBgLogger(BuildLogger(config))
}
