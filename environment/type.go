package environment

type ModeType string

const (
	ModeProduction  ModeType = "production"
	ModeStaging     ModeType = "staging"
	ModeDevelopment ModeType = "development"
)

type LogStyleType string

const (
	LogStyleJSON LogStyleType = "json"
	LogStyleText LogStyleType = "text"
)
