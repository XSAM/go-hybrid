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

// DevelopmentAndJSONConfig set background logger to development mode with JSON style.
func DevelopmentAndJSONConfig() Config {
	config := zap.NewDevelopmentConfig()
	config.Encoding = "json"
	// Keeping the development JSON key naming consistent with production JSON key
	config.EncoderConfig = zap.NewProductionEncoderConfig()

	return Config{ZapConfig: config, ZapLevel: zapcore.DebugLevel}
}

// ProductionAndJSONConfig set background logger to production mode with JSON style.
func ProductionAndJSONConfig() Config {
	config := zap.NewProductionConfig()

	return Config{ZapConfig: config, ZapLevel: zapcore.InfoLevel}
}

// DevelopmentAndTextConfig set background logger to development mode with text style.
func DevelopmentAndTextConfig() Config {
	config := zap.NewDevelopmentConfig()
	if runtime.GOOS != "windows" {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.DisableStacktrace = true

	return Config{ZapConfig: config, ZapLevel: zapcore.DebugLevel}
}

// ProductionAndTextConfig set background logger to production mode with text style.
func ProductionAndTextConfig() Config {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	if runtime.GOOS != "windows" {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.DisableCaller = true
	config.DisableStacktrace = true

	// Disable time key
	config.EncoderConfig.TimeKey = ""

	// Change time and duration encoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

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
