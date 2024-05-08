package config

import (
	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/crypto"
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	"github.com/alexfalkowski/go-service/database/sql"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/feature"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/security/token"
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/telemetry/logger/zap"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

// NewConfig for config.
func NewConfig(i *cmd.InputConfig) (*Config, error) {
	c := &Config{}

	return c, i.Unmarshal(c)
}

// Config for the service.
type Config struct {
	Environment env.Environment   `yaml:"environment,omitempty" json:"environment,omitempty" toml:"environment,omitempty"`
	Debug       *debug.Config     `yaml:"debug,omitempty" json:"debug,omitempty" toml:"debug,omitempty"`
	Cache       *cache.Config     `yaml:"cache,omitempty" json:"cache,omitempty" toml:"cache,omitempty"`
	Crypto      *crypto.Config    `yaml:"crypto,omitempty" json:"crypto,omitempty" toml:"crypto,omitempty"`
	Feature     *feature.Config   `yaml:"feature,omitempty" json:"feature,omitempty" toml:"feature,omitempty"`
	Hooks       *hooks.Config     `yaml:"hooks,omitempty" json:"hooks,omitempty" toml:"hooks,omitempty"`
	Limiter     *limiter.Config   `yaml:"limiter,omitempty" json:"limiter,omitempty" toml:"limiter,omitempty"`
	SQL         *sql.Config       `yaml:"sql,omitempty" json:"sql,omitempty" toml:"sql,omitempty"`
	Telemetry   *telemetry.Config `yaml:"telemetry,omitempty" json:"telemetry,omitempty" toml:"telemetry,omitempty"`
	Time        *time.Config      `yaml:"time,omitempty" json:"time,omitempty" toml:"time,omitempty"`
	Token       *token.Config     `yaml:"token,omitempty" json:"token,omitempty" toml:"token,omitempty"`
	Transport   *transport.Config `yaml:"transport,omitempty" json:"transport,omitempty" toml:"transport,omitempty"`
}

func aesConfig(cfg *Config) *aes.Config {
	if !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.AES
}

func debugConfig(cfg *Config) *debug.Config {
	return cfg.Debug
}

func ed25519Config(cfg *Config) *ed25519.Config {
	if !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.Ed25519
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

func hmacConfig(cfg *Config) *hmac.Config {
	if !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.HMAC
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

func loggerConfig(cfg *Config) *zap.Config {
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

func rsaConfig(cfg *Config) *rsa.Config {
	if !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.RSA
}

func timeConfig(cfg *Config) *time.Config {
	if !time.IsEnabled(cfg.Time) {
		return nil
	}

	return cfg.Time
}

func pgConfig(cfg *Config) *pg.Config {
	if !sql.IsEnabled(cfg.SQL) {
		return nil
	}

	return cfg.SQL.PG
}

func redisConfig(cfg *Config) *redis.Config {
	if !cache.IsEnabled(cfg.Cache) {
		return nil
	}

	return cfg.Cache.Redis
}

func ristrettoConfig(cfg *Config) *ristretto.Config {
	if !cache.IsEnabled(cfg.Cache) {
		return nil
	}

	return cfg.Cache.Ristretto
}

func tokenConfig(cfg *Config) *token.Config {
	return cfg.Token
}

func tracerConfig(cfg *Config) *tracer.Config {
	if !telemetry.IsEnabled(cfg.Telemetry) {
		return nil
	}

	return cfg.Telemetry.Tracer
}
