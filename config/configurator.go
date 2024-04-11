package config

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/security/token"
	"github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

// NewConfigurator for config.
func NewConfigurator(i *cmd.InputConfig) (Configurator, error) {
	c := &Config{}

	return c, i.Unmarshal(c)
}

// Configurator for config.
type Configurator interface {
	EnvironmentConfig() env.Environment
	DebugConfig() *debug.Config
	RedisConfig() *redis.Config
	RistrettoConfig() *ristretto.Config
	FeatureConfig() *feature.Config
	HooksConfig() *hooks.Config
	PGConfig() *pg.Config
	LoggerConfig() *zap.Config
	MetricsConfig() *metrics.Config
	TokenConfig() *token.Config
	TracerConfig() *tracer.Config
	GRPCConfig() *grpc.Config
	HTTPConfig() *http.Config
}

func environmentConfig(cfg Configurator) env.Environment {
	return cfg.EnvironmentConfig()
}

func debugConfig(cfg Configurator) *debug.Config {
	return cfg.DebugConfig()
}

func featureConfig(cfg Configurator) *feature.Config {
	return cfg.FeatureConfig()
}

func hooksConfig(cfg Configurator) *hooks.Config {
	return cfg.HooksConfig()
}

func redisConfig(cfg Configurator) *redis.Config {
	return cfg.RedisConfig()
}

func ristrettoConfig(cfg Configurator) *ristretto.Config {
	return cfg.RistrettoConfig()
}

func pgConfig(cfg Configurator) *pg.Config {
	return cfg.PGConfig()
}

func loggerConfig(cfg Configurator) *zap.Config {
	return cfg.LoggerConfig()
}

func metricsConfig(cfg Configurator) *metrics.Config {
	return cfg.MetricsConfig()
}

func tokenConfig(cfg Configurator) *token.Config {
	return cfg.TokenConfig()
}

func tracerConfig(cfg Configurator) *tracer.Config {
	return cfg.TracerConfig()
}

func grpcConfig(cfg Configurator) *grpc.Config {
	return cfg.GRPCConfig()
}

func httpConfig(cfg Configurator) *http.Config {
	return cfg.HTTPConfig()
}
