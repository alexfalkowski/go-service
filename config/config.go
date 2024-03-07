package config

import (
	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/database/sql"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/security/token"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

// Config for the service.
type Config struct {
	Environment env.Environment  `yaml:"environment" json:"environment" toml:"environment"`
	Debug       debug.Config     `yaml:"debug" json:"debug" toml:"debug"`
	Cache       cache.Config     `yaml:"cache" json:"cache" toml:"cache"`
	Feature     feature.Config   `yaml:"feature" json:"feature" toml:"feature"`
	SQL         sql.Config       `yaml:"sql" json:"sql" toml:"sql"`
	Telemetry   telemetry.Config `yaml:"telemetry" json:"telemetry" toml:"telemetry"`
	Token       token.Config     `yaml:"token" json:"token" toml:"token"`
	Transport   transport.Config `yaml:"transport" json:"transport" toml:"transport"`
}

func (cfg *Config) EnvironmentConfig() env.Environment {
	return cfg.Environment
}

func (cfg *Config) DebugConfig() *debug.Config {
	return &cfg.Debug
}

func (cfg *Config) RedisConfig() *redis.Config {
	return &cfg.Cache.Redis
}

func (cfg *Config) RistrettoConfig() *ristretto.Config {
	return &cfg.Cache.Ristretto
}

func (cfg *Config) PGConfig() *pg.Config {
	return &cfg.SQL.PG
}

func (cfg *Config) FeatureConfig() *feature.Config {
	return &cfg.Feature
}

func (cfg *Config) TracerConfig() *tracer.Config {
	return &cfg.Telemetry.Tracer
}

func (cfg *Config) LoggerConfig() *zap.Config {
	return &cfg.Telemetry.Logger
}

func (cfg *Config) TokenConfig() *token.Config {
	return &cfg.Token
}

func (cfg *Config) GRPCConfig() *grpc.Config {
	return &cfg.Transport.GRPC
}

func (cfg *Config) HTTPConfig() *http.Config {
	return &cfg.Transport.HTTP
}
