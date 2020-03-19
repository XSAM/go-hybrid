package log

import (
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	ZapConfig zap.Config
	ZapLevel  zapcore.Level
}

// DevelopmentConfig set background logger to develop mode.
func DevelopmentConfig() Config {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return Config{ZapConfig: config, ZapLevel: zapcore.DebugLevel}
}

// ProductionConfig set background logger to develop mode.
func ProductionConfig() Config {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return Config{ZapConfig: config, ZapLevel: zapcore.InfoLevel}
}

// CLIToolDevelopmentConfig set background logger to develop mode for CLI.
func CLIToolDevelopmentConfig() Config {
	config := zap.NewDevelopmentConfig()
	if runtime.GOOS != "windows" {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.DisableCaller = true
	config.DisableStacktrace = true

	return Config{ZapConfig: config, ZapLevel: zapcore.DebugLevel}
}

// CLIToolProductionConfig set background logger to production mode for CLI.
func CLIToolProductionConfig() Config {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	if runtime.GOOS != "windows" {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.DisableCaller = true
	config.DisableStacktrace = true

	// Disable time encoder
	config.EncoderConfig.TimeKey = ""

	return Config{ZapConfig: config, ZapLevel: zapcore.InfoLevel}
}

func BuildLogger(config Config) *Core {
	GetLevels().Set(config.ZapLevel)
	// Dynamic log level
	config.ZapConfig.Level = *GetLevels().Get()

	zapLogger, err := config.ZapConfig.Build()
	if err != nil {
		panic("init zap logger: " + err.Error())
	}
	return &Core{Logger: zapLogger}
}
