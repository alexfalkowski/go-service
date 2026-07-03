package config

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

func attributeMap(cfg *Config) attributes.Map {
	if cfg.Telemetry.IsEnabled() {
		return cfg.Telemetry.Attributes
	}
	return nil
}

func loggerConfig(cfg *Config, fs *os.FS) *logger.Config {
	if cfg.Telemetry.IsEnabled() && cfg.Telemetry.Logger.IsEnabled() {
		cfg.Telemetry.Logger.Headers.MustSecrets(fs)
		return cfg.Telemetry.Logger
	}
	return nil
}

func metricsConfig(cfg *Config, fs *os.FS) *metrics.Config {
	if cfg.Telemetry.IsEnabled() && cfg.Telemetry.Metrics.IsEnabled() {
		cfg.Telemetry.Metrics.Headers.MustSecrets(fs)
		return cfg.Telemetry.Metrics
	}
	return nil
}

func propagationConfig(cfg *Config) *telemetry.PropagationConfig {
	if cfg.Telemetry.IsEnabled() {
		return cfg.Telemetry.Propagation
	}
	return nil
}

func tracerConfig(cfg *Config, fs *os.FS) *tracer.Config {
	if cfg.Telemetry.IsEnabled() && cfg.Telemetry.Tracer.IsEnabled() {
		cfg.Telemetry.Tracer.Headers.MustSecrets(fs)
		return cfg.Telemetry.Tracer
	}
	return nil
}
