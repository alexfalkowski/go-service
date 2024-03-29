package zap

import (
	"time"

	"github.com/alexfalkowski/go-service/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// IsEnabled for zap.
func IsEnabled(c *Config) bool {
	return c != nil && c.Enabled
}

// Config for zap.
type Config struct {
	Enabled bool   `yaml:"enabled,omitempty" json:"enabled,omitempty" toml:"enabled,omitempty"`
	Level   string `yaml:"level,omitempty" json:"level,omitempty" toml:"level,omitempty"`
}

// NewConfig for zap.
func NewConfig(env env.Environment, config *Config) (zap.Config, error) {
	if !IsEnabled(config) {
		return zap.Config{}, nil
	}

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
		enc.AppendString(t.Format(time.RFC3339))
	})

	return cfg, nil
}
