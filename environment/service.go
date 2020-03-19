package environment

import (
	"github.com/gin-gonic/gin"

	"github.com/XSAM/go-hybrid/log"
)

var Mode = PRODUCTION_MODE
var Interaction = NORMAL_INTERACTIION

func DevelopmentMode() {
	log.BuildAndSetBgLogger(log.DevelopmentConfig())
	gin.SetMode(gin.DebugMode)
	Mode = DEVELOPMENT_MODE
	Interaction = NORMAL_INTERACTIION
}

func ProductionMode() {
	log.BuildAndSetBgLogger(log.ProductionConfig())
	gin.SetMode(gin.ReleaseMode)
	Mode = PRODUCTION_MODE
	Interaction = NORMAL_INTERACTIION
}

func CLIToolDevelopmentMode() {
	log.BuildAndSetBgLogger(log.CLIToolDevelopmentConfig())
	gin.SetMode(gin.DebugMode)
	Mode = DEVELOPMENT_MODE
	Interaction = CLI_INTERACTION
}

func CLIToolProductionMode() {
	log.BuildAndSetBgLogger(log.CLIToolProductionConfig())
	gin.SetMode(gin.ReleaseMode)
	Mode = PRODUCTION_MODE
	Interaction = CLI_INTERACTION
}
