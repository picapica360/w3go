package zap

import (
	"os"

	"go.uber.org/zap/zapcore"
)

// NewConsoleAdapter create a console adapter.
func NewConsoleAdapter() func() zapcore.Core {
	return func() zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(DefaultEncoderConfig),
			zapcore.AddSync(os.Stdout),
			DefaultLevelEnablerFunc(),
		)
	}
}
