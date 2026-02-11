package config

import (
	cache "github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/crypto"
	"github.com/alexfalkowski/go-service/v2/database/sql"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/feature"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/types/ptr"
	"github.com/alexfalkowski/go-service/v2/types/structs"
)

// NewConfig will decode and check its validity.
func NewConfig[T comparable](decoder Decoder, validator *Validator) (*T, error) {
	config := ptr.Zero[T]()
	if err := decoder.Decode(config); err != nil {
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

// Config is the top-level configuration for a go-service based service.
//
// It composes feature configurations from other packages (debug, cache, crypto, etc.).
// All pointer fields are optional; when a sub-config is nil it is generally treated as disabled
// (see each sub-config's IsEnabled method, where applicable).
type Config struct {
	// Debug configures the debug server and related debugging endpoints.
	Debug *debug.Config `yaml:"debug,omitempty" json:"debug,omitempty" toml:"debug,omitempty"`

	// Cache configures the cache subsystem (implementation kind and implementation-specific options).
	Cache *cache.Config `yaml:"cache,omitempty" json:"cache,omitempty" toml:"cache,omitempty"`

	// Crypto configures cryptographic primitives used by the service (for example HMAC, RSA, Ed25519, SSH, AES).
	Crypto *crypto.Config `yaml:"crypto,omitempty" json:"crypto,omitempty" toml:"crypto,omitempty"`

	// Feature configures feature client behavior used by some internal feature integrations.
	Feature *feature.Config `yaml:"feature,omitempty" json:"feature,omitempty" toml:"feature,omitempty"`

	// Hooks configures webhook behavior (for example the shared secret used to validate incoming hooks).
	Hooks *hooks.Config `yaml:"hooks,omitempty" json:"hooks,omitempty" toml:"hooks,omitempty"`

	// ID configures ID generation (for example which generator kind to use).
	ID *id.Config `yaml:"id,omitempty" json:"id,omitempty" toml:"id,omitempty"`

	// SQL configures SQL database access (for example PostgreSQL connection pools and DSNs).
	SQL *sql.Config `yaml:"sql,omitempty" json:"sql,omitempty" toml:"sql,omitempty"`

	// Telemetry configures logging, tracing, and metrics.
	Telemetry *telemetry.Config `yaml:"telemetry,omitempty" json:"telemetry,omitempty" toml:"telemetry,omitempty"`

	// Time configures a network time provider (for example NTP/NTS) used by the service.
	Time *time.Config `yaml:"time,omitempty" json:"time,omitempty" toml:"time,omitempty"`

	// Transport configures inbound/outbound transports (HTTP and gRPC).
	Transport *transport.Config `yaml:"transport,omitempty" json:"transport,omitempty" toml:"transport,omitempty"`

	// Environment is the service environment (for example local/dev/stage/prod) used to drive environment-specific behavior.
	Environment env.Environment `yaml:"environment,omitempty" json:"environment,omitempty" toml:"environment,omitempty"`
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
	if cfg.Transport.IsEnabled() {
		return cfg.Transport.GRPC
	}
	return nil
}

func idConfig(cfg *Config) *id.Config {
	return cfg.ID
}

func hooksConfig(cfg *Config) *hooks.Config {
	return cfg.Hooks
}

func httpConfig(cfg *Config) *http.Config {
	if cfg.Transport.IsEnabled() {
		return cfg.Transport.HTTP
	}
	return nil
}

func pgConfig(cfg *Config) *pg.Config {
	if cfg.SQL.IsEnabled() {
		return cfg.SQL.PG
	}
	return nil
}

func timeConfig(cfg *Config) *time.Config {
	if cfg.Time.IsEnabled() {
		return cfg.Time
	}
	return nil
}
