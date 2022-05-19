package config

import (
	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/database/sql"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/security/auth0"
	"github.com/alexfalkowski/go-service/trace"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/nsq"
)

// Config for the service.
type Config struct {
	Runtime   string           `yaml:"runtime"`
	Cache     cache.Config     `yaml:"cache"`
	Security  security.Config  `yaml:"security"`
	SQL       sql.Config       `yaml:"sql"`
	Trace     trace.Config     `yaml:"trace"`
	Transport transport.Config `yaml:"transport"`
}

// GetRuntime of the application. The values can be os or container.
// This is used to determine how we watch the configuration.
func (cfg *Config) GetRuntime() string {
	runtime := cfg.Runtime
	if runtime == "" {
		return "os"
	}

	return runtime
}

func (cfg *Config) RedisConfig() *redis.Config {
	return &cfg.Cache.Redis
}

func (cfg *Config) RistrettoConfig() *ristretto.Config {
	return &cfg.Cache.Ristretto
}

func (cfg *Config) Auth0Config() *auth0.Config {
	return &cfg.Security.Auth0
}

func (cfg *Config) PGConfig() *pg.Config {
	return &cfg.SQL.PG
}

func (cfg *Config) OpentracingConfig() *opentracing.Config {
	return &cfg.Trace.Opentracing
}

func (cfg *Config) TransportConfig() *transport.Config {
	return &cfg.Transport
}

func (cfg *Config) GRPCConfig() *grpc.Config {
	return &cfg.Transport.GRPC
}

func (cfg *Config) HTTPConfig() *http.Config {
	return &cfg.Transport.HTTP
}

func (cfg *Config) NSQConfig() *nsq.Config {
	return &cfg.Transport.NSQ
}
