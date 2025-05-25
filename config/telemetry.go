package config

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

func loggerConfig(cfg *Config, fs *os.FS) *logger.Config {
	if !telemetry.IsEnabled(cfg.Telemetry) || !logger.IsEnabled(cfg.Telemetry.Logger) {
		return nil
	}

	err := cfg.Telemetry.Logger.Headers.Secrets(fs)
	runtime.Must(err)

	return cfg.Telemetry.Logger
}

func metricsConfig(cfg *Config, fs *os.FS) *metrics.Config {
	if !telemetry.IsEnabled(cfg.Telemetry) || !metrics.IsEnabled(cfg.Telemetry.Metrics) {
		return nil
	}

	err := cfg.Telemetry.Metrics.Headers.Secrets(fs)
	runtime.Must(err)

	return cfg.Telemetry.Metrics
}

func tracerConfig(cfg *Config, fs *os.FS) *tracer.Config {
	if !telemetry.IsEnabled(cfg.Telemetry) || !tracer.IsEnabled(cfg.Telemetry.Tracer) {
		return nil
	}

	err := cfg.Telemetry.Tracer.Headers.Secrets(fs)
	runtime.Must(err)

	return cfg.Telemetry.Tracer
}
