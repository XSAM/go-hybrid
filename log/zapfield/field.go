package zapfield

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/XSAM/go-hybrid/builtinutil"
	"github.com/XSAM/go-hybrid/environment"
	"github.com/XSAM/go-hybrid/errorw"
)

func Stack() zap.Field {
	return zap.String("stack", string(builtinutil.Stack(2)))
}

// Error add error field
func Error(err error) zap.Field {
	if environment.Interaction == environment.CLI_INTERACTION {
		// CLI
		switch environment.Mode {
		case environment.DEVELOPMENT_MODE:
			// Normal encoding
			switch e := err.(type) {
			case *errorw.Error:
				return zap.Field{Key: "error", Type: zapcore.ObjectMarshalerType, Interface: e}
			}
		case environment.PRODUCTION_MODE:
			return zap.Field{Key: "error", Type: zapcore.StringType, String: err.Error()}
		}
	} else {
		// Normal
		switch e := err.(type) {
		case *errorw.Error:
			return zap.Field{Key: "error", Type: zapcore.ObjectMarshalerType, Interface: e}
		}
	}
	return zap.Error(err)
}
