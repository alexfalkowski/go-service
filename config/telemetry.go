package config

import (
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
)

func (cfg *Config) LoggerConfig() *zap.Config {
	if !telemetry.IsEnabled(cfg.Telemetry) {
		return nil
	}

	return cfg.Telemetry.Logger
}

func (cfg *Config) MetricsConfig() *metrics.Config {
	if !telemetry.IsEnabled(cfg.Telemetry) {
		return nil
	}

	return cfg.Telemetry.Metrics
}

func (cfg *Config) TracerConfig() *tracer.Config {
	if !telemetry.IsEnabled(cfg.Telemetry) {
		return nil
	}

	return cfg.Telemetry.Tracer
}
