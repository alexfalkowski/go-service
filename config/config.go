package config

import (
	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/database/sql"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/security/token"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

// Config for the service.
type Config struct {
	Environment env.Environment   `yaml:"environment,omitempty" json:"environment,omitempty" toml:"environment,omitempty"`
	Debug       *debug.Config     `yaml:"debug,omitempty" json:"debug,omitempty" toml:"debug,omitempty"`
	Cache       *cache.Config     `yaml:"cache,omitempty" json:"cache,omitempty" toml:"cache,omitempty"`
	Feature     *feature.Config   `yaml:"feature,omitempty" json:"feature,omitempty" toml:"feature,omitempty"`
	Hooks       *hooks.Config     `yaml:"hooks,omitempty" json:"hooks,omitempty" toml:"hooks,omitempty"`
	Limiter     *limiter.Config   `yaml:"limiter,omitempty" json:"limiter,omitempty" toml:"limiter,omitempty"`
	SQL         *sql.Config       `yaml:"sql,omitempty" json:"sql,omitempty" toml:"sql,omitempty"`
	Telemetry   *telemetry.Config `yaml:"telemetry,omitempty" json:"telemetry,omitempty" toml:"telemetry,omitempty"`
	Token       *token.Config     `yaml:"token,omitempty" json:"token,omitempty" toml:"token,omitempty"`
	Transport   *transport.Config `yaml:"transport,omitempty" json:"transport,omitempty" toml:"transport,omitempty"`
}

func (cfg *Config) DebugConfig() *debug.Config {
	return cfg.Debug
}

func (cfg *Config) EnvironmentConfig() env.Environment {
	return cfg.Environment
}

func (cfg *Config) FeatureConfig() *feature.Config {
	return cfg.Feature
}

func (cfg *Config) GRPCConfig() *grpc.Config {
	if cfg.Transport == nil {
		return nil
	}

	return cfg.Transport.GRPC
}

func (cfg *Config) HooksConfig() *hooks.Config {
	return cfg.Hooks
}

func (cfg *Config) HTTPConfig() *http.Config {
	if cfg.Transport == nil {
		return nil
	}

	return cfg.Transport.HTTP
}

func (cfg *Config) LimiterConfig() *limiter.Config {
	return cfg.Limiter
}

func (cfg *Config) PGConfig() *pg.Config {
	if cfg.SQL == nil {
		return nil
	}

	return cfg.SQL.PG
}

func (cfg *Config) TokenConfig() *token.Config {
	return cfg.Token
}
