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
	"github.com/alexfalkowski/go-service/limiter"
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
	DebugConfig() *debug.Config
	EnvironmentConfig() env.Environment
	FeatureConfig() *feature.Config
	GRPCConfig() *grpc.Config
	HooksConfig() *hooks.Config
	HTTPConfig() *http.Config
	LimiterConfig() *limiter.Config
	LoggerConfig() *zap.Config
	MetricsConfig() *metrics.Config
	PGConfig() *pg.Config
	RedisConfig() *redis.Config
	RistrettoConfig() *ristretto.Config
	TokenConfig() *token.Config
	TracerConfig() *tracer.Config
}

func debugConfig(cfg Configurator) *debug.Config {
	return cfg.DebugConfig()
}

func environmentConfig(cfg Configurator) env.Environment {
	return cfg.EnvironmentConfig()
}

func featureConfig(cfg Configurator) *feature.Config {
	return cfg.FeatureConfig()
}

func grpcConfig(cfg Configurator) *grpc.Config {
	return cfg.GRPCConfig()
}

func hooksConfig(cfg Configurator) *hooks.Config {
	return cfg.HooksConfig()
}

func httpConfig(cfg Configurator) *http.Config {
	return cfg.HTTPConfig()
}

func limiterConfig(cfg Configurator) *limiter.Config {
	return cfg.LimiterConfig()
}

func loggerConfig(cfg Configurator) *zap.Config {
	return cfg.LoggerConfig()
}

func metricsConfig(cfg Configurator) *metrics.Config {
	return cfg.MetricsConfig()
}

func pgConfig(cfg Configurator) *pg.Config {
	return cfg.PGConfig()
}

func redisConfig(cfg Configurator) *redis.Config {
	return cfg.RedisConfig()
}

func ristrettoConfig(cfg Configurator) *ristretto.Config {
	return cfg.RistrettoConfig()
}

func tokenConfig(cfg Configurator) *token.Config {
	return cfg.TokenConfig()
}

func tracerConfig(cfg Configurator) *tracer.Config {
	return cfg.TracerConfig()
}
