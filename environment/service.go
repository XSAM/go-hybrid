package environment

import (
	"github.com/gin-gonic/gin"

	"github.com/XSAM/go-hybrid/log"
)

var Mode = ModeProduction
var LogStyle = LogStyleJSON

func DevelopmentMode() {
	gin.SetMode(gin.DebugMode)
	Mode = ModeDevelopment
}

func ProductionMode() {
	gin.SetMode(gin.ReleaseMode)
	Mode = ModeProduction
}

func JSONLogStyle() {
	LogStyle = LogStyleJSON
	switch Mode {
	case ModeDevelopment:
		log.BuildAndSetBgLogger(log.DevelopmentAndJSONConfig())
	case ModeProduction, ModeStaging:
		log.BuildAndSetBgLogger(log.ProductionAndJSONConfig())
	}
}

func TextLogStyle() {
	LogStyle = LogStyleText
	switch Mode {
	case ModeDevelopment:
		log.BuildAndSetBgLogger(log.DevelopmentAndTextConfig())
	case ModeProduction, ModeStaging:
		log.BuildAndSetBgLogger(log.ProductionAndTextConfig())
	}
}
