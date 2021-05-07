package zap

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewConfig for zap.
func NewConfig() zap.Config {
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableCaller = true
	zapConfig.DisableStacktrace = true
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(time.RFC3339))
	})

	return zapConfig
}
