package environment

type ModeType string

const (
	PRODUCTION_MODE  ModeType = "production"
	STAGING_MODE              = "staging"
	DEVELOPMENT_MODE          = "development"
)

type InteractionType string

const (
	NORMAL_INTERACTIION InteractionType = "normal"
	CLI_INTERACTION                     = "CLI"
)
