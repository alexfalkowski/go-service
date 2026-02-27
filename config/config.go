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

// NewConfig decodes configuration into a newly allocated *T and validates it.
//
// Decoding is performed using the provided Decoder, which typically resolves configuration from:
//   - a file (via "file:<path>"),
//   - an environment variable (via "env:<ENV_VAR>"), or
//   - the default lookup (see the config package documentation).
//
// After decoding, NewConfig enforces two additional invariants:
//
//  1. Empty detection: if the decoded value is considered empty (see structs.IsEmpty), NewConfig returns
//     ErrInvalidConfig. This guards against accidentally starting with a zero-value configuration when the
//     input is missing or does not populate any fields.
//
//  2. Validation: the decoded value is validated using the provided Validator (go-playground/validator).
//     Any validation errors are returned to the caller.
//
// On success, NewConfig returns the validated configuration value.
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

// Config is the standard top-level configuration shape for a go-service based service.
//
// It composes feature configurations from other packages (debug, cache, crypto, telemetry, transports, etc.)
// so services can embed a single `config.Config` and get consistent wiring.
//
// # Optional pointers and "enabled" semantics
//
// Most fields are pointers and are intentionally optional. A nil sub-config is generally treated as
// "disabled" by the corresponding subsystem. Many sub-config types expose an `IsEnabled()` method
// whose convention is `return c != nil`.
//
// # Transport-derived projections
//
// Some subsystems are nested (for example HTTP/GRPC under Transport, and PG under SQL). The config
// module wiring provides small projection constructors (see config/module.go) that safely return nil
// when the parent sub-config is disabled.
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
