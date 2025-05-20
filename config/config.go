package config

import (
	"errors"

	cache "github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/cli"
	"github.com/alexfalkowski/go-service/crypto"
	"github.com/alexfalkowski/go-service/database/sql"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/types/ptr"
	"github.com/alexfalkowski/go-service/types/structs"
	"github.com/alexfalkowski/go-service/types/validator"
)

// ErrInvalidConfig when decoding fails.
var ErrInvalidConfig = errors.New("config: invalid format")

// NewConfig will decode and check its validity.
func NewConfig[T comparable](input *cli.InputConfig, validator *validator.Validator) (*T, error) {
	config := ptr.Zero[T]()
	if err := input.Decode(config); err != nil {
		return nil, err
	}

	if structs.IsEmpty(config) {
		return nil, ErrInvalidConfig
	}

	if err := validator.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}

// Config for the service.
type Config struct {
	Debug       *debug.Config     `yaml:"debug,omitempty" json:"debug,omitempty" toml:"debug,omitempty"`
	Cache       *cache.Config     `yaml:"cache,omitempty" json:"cache,omitempty" toml:"cache,omitempty"`
	Crypto      *crypto.Config    `yaml:"crypto,omitempty" json:"crypto,omitempty" toml:"crypto,omitempty"`
	Feature     *feature.Config   `yaml:"feature,omitempty" json:"feature,omitempty" toml:"feature,omitempty"`
	Hooks       *hooks.Config     `yaml:"hooks,omitempty" json:"hooks,omitempty" toml:"hooks,omitempty"`
	ID          *id.Config        `yaml:"id,omitempty" json:"id,omitempty" toml:"id,omitempty"`
	Limiter     *limiter.Config   `yaml:"limiter,omitempty" json:"limiter,omitempty" toml:"limiter,omitempty"`
	SQL         *sql.Config       `yaml:"sql,omitempty" json:"sql,omitempty" toml:"sql,omitempty"`
	Telemetry   *telemetry.Config `yaml:"telemetry,omitempty" json:"telemetry,omitempty" toml:"telemetry,omitempty"`
	Time        *time.Config      `yaml:"time,omitempty" json:"time,omitempty" toml:"time,omitempty"`
	Token       *token.Config     `yaml:"token,omitempty" json:"token,omitempty" toml:"token,omitempty"`
	Transport   *transport.Config `yaml:"transport,omitempty" json:"transport,omitempty" toml:"transport,omitempty"`
	Environment env.Environment   `yaml:"environment,omitempty" json:"environment,omitempty" toml:"environment,omitempty"`
}

func cacheConfig(cfg *Config) *cache.Config {
	return cfg.Cache
}

func debugConfig(cfg *Config) *debug.Config {
	return cfg.Debug
}

func environmentConfig(cfg *Config) env.Environment {
	return cfg.Environment
}

func featureConfig(cfg *Config) *feature.Config {
	return cfg.Feature
}

func grpcConfig(cfg *Config) *grpc.Config {
	if !transport.IsEnabled(cfg.Transport) {
		return nil
	}

	return cfg.Transport.GRPC
}

func idConfig(cfg *Config) *id.Config {
	return cfg.ID
}

func hooksConfig(cfg *Config) *hooks.Config {
	return cfg.Hooks
}

func httpConfig(cfg *Config) *http.Config {
	if !transport.IsEnabled(cfg.Transport) {
		return nil
	}

	return cfg.Transport.HTTP
}

func limiterConfig(cfg *Config) *limiter.Config {
	return cfg.Limiter
}

func loggerConfig(cfg *Config) *logger.Config {
	if !telemetry.IsEnabled(cfg.Telemetry) {
		return nil
	}

	return cfg.Telemetry.Logger
}

func metricsConfig(cfg *Config) *metrics.Config {
	if !telemetry.IsEnabled(cfg.Telemetry) {
		return nil
	}

	return cfg.Telemetry.Metrics
}

func pgConfig(cfg *Config) *pg.Config {
	if !sql.IsEnabled(cfg.SQL) {
		return nil
	}

	return cfg.SQL.PG
}

func timeConfig(cfg *Config) *time.Config {
	if !time.IsEnabled(cfg.Time) {
		return nil
	}

	return cfg.Time
}

func tracerConfig(cfg *Config) *tracer.Config {
	if !telemetry.IsEnabled(cfg.Telemetry) {
		return nil
	}

	return cfg.Telemetry.Tracer
}
