package telemetry

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/version"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerParams for telemetry.
type LoggerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Version   version.Version
}

// NewLogger for telemetry.
func NewLogger(params LoggerParams) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(time.RFC3339))
	})

	logger, err := cfg.Build(zap.Fields(zap.String("name", os.ExecutableName()), zap.String("version", string(params.Version))))
	if err != nil {
		return nil, err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = logger.Sync()

			return nil
		},
	})

	return logger, nil
}
