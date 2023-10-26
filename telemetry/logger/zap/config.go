package zap

import (
	"time"

	"github.com/alexfalkowski/go-service/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config for zap.
type Config struct {
	Level string `yaml:"level" json:"level" toml:"level"`
}

// NewConfig for zap.
func NewConfig(env env.Environment, config *Config) (zap.Config, error) {
	l, err := zap.ParseAtomicLevel(config.Level)
	if err != nil {
		return zap.Config{}, err
	}

	var cfg zap.Config

	if env.IsDevelopment() {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.Level = l
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(time.RFC3339))
	})

	return cfg, nil
}
