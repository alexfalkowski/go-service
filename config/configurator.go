package config

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/security/auth0"
	"github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/nsq"
)

// NewConfigurator for config.
func NewConfigurator(i *cmd.InputConfig) (Configurator, error) {
	c := &Config{}

	if err := i.Unmarshal(c); err != nil {
		return nil, err
	}

	return c, nil
}

// Configurator for config.
type Configurator interface {
	RedisConfig() *redis.Config
	RistrettoConfig() *ristretto.Config
	Auth0Config() *auth0.Config
	PGConfig() *pg.Config
	LoggerConfig() *zap.Config
	TracerConfig() *tracer.Config
	TransportConfig() *transport.Config
	GRPCConfig() *grpc.Config
	HTTPConfig() *http.Config
	NSQConfig() *nsq.Config
}

func redisConfig(cfg Configurator) *redis.Config {
	return cfg.RedisConfig()
}

func ristrettoConfig(cfg Configurator) *ristretto.Config {
	return cfg.RistrettoConfig()
}

func auth0Config(cfg Configurator) *auth0.Config {
	return cfg.Auth0Config()
}

func pgConfig(cfg Configurator) *pg.Config {
	return cfg.PGConfig()
}

func loggerConfig(cfg Configurator) *zap.Config {
	return cfg.LoggerConfig()
}

func tracerConfig(cfg Configurator) *tracer.Config {
	return cfg.TracerConfig()
}

func transportConfig(cfg Configurator) *transport.Config {
	return cfg.TransportConfig()
}

func grpcConfig(cfg Configurator) *grpc.Config {
	return cfg.GRPCConfig()
}

func httpConfig(cfg Configurator) *http.Config {
	return cfg.HTTPConfig()
}

func nsqConfig(cfg Configurator) *nsq.Config {
	return cfg.NSQConfig()
}
