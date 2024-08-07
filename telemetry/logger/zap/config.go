package zap

import (
	"time"

	"github.com/alexfalkowski/go-service/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// IsEnabled for zap.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for zap.
type Config struct {
	Level string `yaml:"level,omitempty" json:"level,omitempty" toml:"level,omitempty"`
}

// NewConfig for zap. If disabled returns nil and ignored by logger.
//
//nolint:nilnil
func NewConfig(env env.Environment, config *Config) (*zap.Config, error) {
	if !IsEnabled(config) {
		return nil, nil
	}

	l, err := zap.ParseAtomicLevel(config.Level)
	if err != nil {
		return nil, err
	}

	var cfg zap.Config

	if env.IsDevelopment() {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.Level = l
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	})

	return &cfg, nil
}
