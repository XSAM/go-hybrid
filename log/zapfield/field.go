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
	if err == nil {
		return zap.Skip()
	}
	
	if environment.LogStyle == environment.LogStyleText {
		// Text style
		switch environment.Mode {
		case environment.ModeDevelopment:
			switch e := err.(type) {
			case *errorw.Error:
				return zap.Field{Key: "error", Type: zapcore.ObjectMarshalerType, Interface: e}
			}
		case environment.ModeProduction, environment.ModeStaging:
			return zap.Field{Key: "error", Type: zapcore.StringType, String: err.Error()}
		}
	} else {
		// JSON style
		switch e := err.(type) {
		case *errorw.Error:
			return zap.Field{Key: "error", Type: zapcore.ObjectMarshalerType, Interface: e}
		}
	}
	return zap.Error(err)
}
