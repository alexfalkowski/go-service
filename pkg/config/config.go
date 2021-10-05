package config

import (
	"errors"
	"os"

	"github.com/alexfalkowski/go-service/pkg/cache"
	"github.com/alexfalkowski/go-service/pkg/cache/redis"
	"github.com/alexfalkowski/go-service/pkg/cache/ristretto"
	"github.com/alexfalkowski/go-service/pkg/security"
	"github.com/alexfalkowski/go-service/pkg/security/auth0"
	"github.com/alexfalkowski/go-service/pkg/sql"
	"github.com/alexfalkowski/go-service/pkg/sql/pg"
	"github.com/alexfalkowski/go-service/pkg/trace"
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/jaeger"
	"github.com/alexfalkowski/go-service/pkg/transport"
	"github.com/alexfalkowski/go-service/pkg/transport/grpc"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq"
	"gopkg.in/yaml.v3"
)

var (
	// ErrMissingConfigFile for config.
	ErrMissingConfigFile = errors.New("missing config file")
)

// Config for the service.
type Config struct {
	Cache     cache.Config     `yaml:"cache"`
	Security  security.Config  `yaml:"security"`
	SQL       sql.Config       `yaml:"sql"`
	Trace     trace.Config     `yaml:"trace"`
	Transport transport.Config `yaml:"transport"`
}

// New config for the service.
func New() (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		return nil, ErrMissingConfigFile
	}

	bytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	cfg := Config{}

	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func redisConfig(cfg *Config) *redis.Config {
	return &cfg.Cache.Redis
}

func ristrettoConfig(cfg *Config) *ristretto.Config {
	return &cfg.Cache.Ristretto
}

func auth0Config(cfg *Config) *auth0.Config {
	return &cfg.Security.Auth0
}

func pgConfig(cfg *Config) *pg.Config {
	return &cfg.SQL.PG
}

func datadogConfig(cfg *Config) *datadog.Config {
	return &cfg.Trace.Opentracing.Datadog
}

func jaegerConfig(cfg *Config) *jaeger.Config {
	return &cfg.Trace.Opentracing.Jaeger
}

func grpcConfig(cfg *Config) *grpc.Config {
	return &cfg.Transport.GRPC
}

func httpConfig(cfg *Config) *http.Config {
	return &cfg.Transport.HTTP
}

func nsqConfig(cfg *Config) *nsq.Config {
	return &cfg.Transport.NSQ
}
