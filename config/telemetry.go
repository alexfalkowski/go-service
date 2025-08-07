package config

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

func loggerConfig(cfg *Config, fs *os.FS) *logger.Config {
	if !cfg.Telemetry.IsEnabled() || !cfg.Telemetry.Logger.IsEnabled() {
		return nil
	}

	err := cfg.Telemetry.Logger.Headers.Secrets(fs)
	runtime.Must(err)
	return cfg.Telemetry.Logger
}

func metricsConfig(cfg *Config, fs *os.FS) *metrics.Config {
	if !cfg.Telemetry.IsEnabled() || !cfg.Telemetry.Metrics.IsEnabled() {
		return nil
	}

	err := cfg.Telemetry.Metrics.Headers.Secrets(fs)
	runtime.Must(err)
	return cfg.Telemetry.Metrics
}

func tracerConfig(cfg *Config, fs *os.FS) *tracer.Config {
	if !cfg.Telemetry.IsEnabled() || !cfg.Telemetry.Tracer.IsEnabled() {
		return nil
	}

	err := cfg.Telemetry.Tracer.Headers.Secrets(fs)
	runtime.Must(err)
	return cfg.Telemetry.Tracer
}
