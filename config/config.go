package config

import (
	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/database/sql"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/security/auth0"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/nsq"
)

// Config for the service.
type Config struct {
	Cache     cache.Config     `yaml:"cache" json:"cache" toml:"cache"`
	Security  security.Config  `yaml:"security" json:"security" toml:"security"`
	SQL       sql.Config       `yaml:"sql" json:"sql" toml:"sql"`
	OTEL      otel.Config      `yaml:"otel" json:"otel" toml:"otel"`
	Transport transport.Config `yaml:"transport" json:"transport" toml:"transport"`
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

func (cfg *Config) OTELConfig() *otel.Config {
	return &cfg.OTEL
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
