package environment

type ModeType string

const (
	ModeProduction  ModeType = "production"
	ModeStaging              = "staging"
	ModeDevelopment          = "development"
)

type LogStyleType string

const (
	LogStyleJSON LogStyleType = "json"
	LogStyleText              = "text"
)
