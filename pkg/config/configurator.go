package config

import (
	"errors"
	"os"

	"github.com/alexfalkowski/go-service/pkg/cache/redis"
	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
	"github.com/alexfalkowski/go-service/pkg/security/auth0"
	"github.com/alexfalkowski/go-service/pkg/sql/pg"
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/jaeger"
	"github.com/alexfalkowski/go-service/pkg/transport/grpc"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq"
)

var (
	// ErrMissingConfigFile for config.
	ErrMissingConfigFile = errors.New("missing config file")
)

// Configurator for config.
type Configurator interface {
	Unmarshal(in []byte) error
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

// NewConfigurator for config.
func NewConfigurator() Configurator {
	return &Config{}
}

// Unmarshal the config.
func Unmarshal(cfg Configurator) error {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		return ErrMissingConfigFile
	}

	bytes, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = cfg.Unmarshal(bytes)
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
