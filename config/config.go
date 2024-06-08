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
	"github.com/alexfalkowski/go-service/crypto/ssh"
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

// IsEnabled for config.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for the service.
type Config struct {
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
	Environment env.Environment   `yaml:"environment,omitempty" json:"environment,omitempty" toml:"environment,omitempty"`
}

func aesConfig(cfg *Config) *aes.Config {
	if !IsEnabled(cfg) || !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.AES
}

func debugConfig(cfg *Config) *debug.Config {
	if !IsEnabled(cfg) {
		return nil
	}

	return cfg.Debug
}

func ed25519Config(cfg *Config) *ed25519.Config {
	if !IsEnabled(cfg) || !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.Ed25519
}

func environmentConfig(cfg *Config) env.Environment {
	if !IsEnabled(cfg) {
		return env.Development
	}

	return cfg.Environment
}

func featureConfig(cfg *Config) *feature.Config {
	if !IsEnabled(cfg) {
		return nil
	}

	return cfg.Feature
}

func grpcConfig(cfg *Config) *grpc.Config {
	if !IsEnabled(cfg) || !transport.IsEnabled(cfg.Transport) {
		return nil
	}

	return cfg.Transport.GRPC
}

func hmacConfig(cfg *Config) *hmac.Config {
	if !IsEnabled(cfg) || !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.HMAC
}

func hooksConfig(cfg *Config) *hooks.Config {
	if !IsEnabled(cfg) {
		return nil
	}

	return cfg.Hooks
}

func httpConfig(cfg *Config) *http.Config {
	if !IsEnabled(cfg) || !transport.IsEnabled(cfg.Transport) {
		return nil
	}

	return cfg.Transport.HTTP
}

func limiterConfig(cfg *Config) *limiter.Config {
	if !IsEnabled(cfg) {
		return nil
	}

	return cfg.Limiter
}

func loggerConfig(cfg *Config) *zap.Config {
	if !IsEnabled(cfg) || !telemetry.IsEnabled(cfg.Telemetry) {
		return nil
	}

	return cfg.Telemetry.Logger
}

func metricsConfig(cfg *Config) *metrics.Config {
	if !IsEnabled(cfg) || !telemetry.IsEnabled(cfg.Telemetry) {
		return nil
	}

	return cfg.Telemetry.Metrics
}

func rsaConfig(cfg *Config) *rsa.Config {
	if !IsEnabled(cfg) || !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.RSA
}

func pgConfig(cfg *Config) *pg.Config {
	if !IsEnabled(cfg) || !sql.IsEnabled(cfg.SQL) {
		return nil
	}

	return cfg.SQL.PG
}

func redisConfig(cfg *Config) *redis.Config {
	if !IsEnabled(cfg) || !cache.IsEnabled(cfg.Cache) {
		return nil
	}

	return cfg.Cache.Redis
}

func ristrettoConfig(cfg *Config) *ristretto.Config {
	if !IsEnabled(cfg) || !cache.IsEnabled(cfg.Cache) {
		return nil
	}

	return cfg.Cache.Ristretto
}

func sshConfig(cfg *Config) *ssh.Config {
	if !IsEnabled(cfg) || !crypto.IsEnabled(cfg.Crypto) {
		return nil
	}

	return cfg.Crypto.SSH
}

func timeConfig(cfg *Config) *time.Config {
	if !IsEnabled(cfg) || !time.IsEnabled(cfg.Time) {
		return nil
	}

	return cfg.Time
}

func tokenConfig(cfg *Config) *token.Config {
	if !IsEnabled(cfg) {
		return nil
	}

	return cfg.Token
}

func tracerConfig(cfg *Config) *tracer.Config {
	if !IsEnabled(cfg) || !telemetry.IsEnabled(cfg.Telemetry) {
		return nil
	}

	return cfg.Telemetry.Tracer
}
