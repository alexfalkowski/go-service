package logger

import (
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/time"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// IsEnabled for logger.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for logger.
type Config struct {
	Level string `yaml:"level,omitempty" json:"level,omitempty" toml:"level,omitempty"`
}

// NewConfig for logger. If disabled returns nil and ignored by logger.
func NewConfig(env env.Environment, config *Config) (*zap.Config, error) {
	if !IsEnabled(config) {
		return nil, nil
	}

	level, err := zap.ParseAtomicLevel(config.Level)
	if err != nil {
		return nil, err
	}

	var cfg zap.Config

	if env.IsDevelopment() {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.DisableStacktrace = true
	cfg.DisableCaller = true
	cfg.Level = level
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	})

	return &cfg, nil
}
