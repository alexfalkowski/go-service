package config

import (
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/security/auth0"
	"github.com/alexfalkowski/go-service/sql/pg"
	"github.com/alexfalkowski/go-service/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/nsq"
	"gopkg.in/yaml.v3"
)

// NewConfigurator for config.
// nolint:ireturn
func NewConfigurator() Configurator {
	return &Config{}
}

// Configurator for config.
type Configurator interface {
	RedisConfig() *redis.Config
	RistrettoConfig() *ristretto.Config
	Auth0Config() *auth0.Config
	PGConfig() *pg.Config
	DatadogConfig() *datadog.Config
	JaegerConfig() *jaeger.Config
	GRPCConfig() *grpc.Config
	HTTPConfig() *http.Config
	NSQConfig() *nsq.Config
}

// UnmarshalFromFile to config.
func UnmarshalFromFile(cfg Configurator) error {
	bytes, err := ReadFile()
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		return err
	}

	return nil
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

func datadogConfig(cfg Configurator) *datadog.Config {
	return cfg.DatadogConfig()
}

func jaegerConfig(cfg Configurator) *jaeger.Config {
	return cfg.JaegerConfig()
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
