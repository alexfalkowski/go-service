package logger

import (
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/alexfalkowski/go-service/v2/telemetry/internal/otlp"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Config configures structured logging and optional log exporting.
//
// This package is responsible for constructing the logger based on Config.
// However, note that this type intentionally does not resolve header secrets by
// itself. If Headers contains go-service "source strings" (for example "env:NAME"
// or "file:/path"), callers are expected to resolve them before constructing
// exporters/handlers (for example via [header.Map.Secrets] or [header.Map.MustSecrets]).
//
// Enablement is modeled by presence: a nil *[Config] means logging is disabled.
type Config struct {
	// Headers contains exporter/request headers.
	//
	// These headers are primarily used by exporter-backed logger kinds (for example
	// "otlp") to pass authentication and/or routing metadata to a collector.
	//
	// Values may be configured as go-service "source strings" (for example "env:NAME",
	// "file:/path", or a literal value). Resolution is performed by the consumer that
	// prepares configuration for use by the logger/exporter.
	Headers header.Map `yaml:"headers,omitempty" json:"headers,omitempty" toml:"headers,omitempty"`

	// TLS configures OTLP/gRPC client transport security.
	//
	// This only applies when Kind is "otlp" and Protocol is "grpc". A non-nil
	// config enables TLS; Cert, Key, and CA values use go-service source strings.
	TLS *tls.Config `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`

	// Kind selects the logger implementation.
	//
	// Supported kinds depend on what the service links in, but this package typically
	// supports:
	//
	//   - "otlp": export logs via OpenTelemetry OTLP (and bridge slog records to OTel).
	//   - "json": write JSON logs to stdout.
	//   - "text": write text logs to stdout.
	//   - "tint": write colorized text logs to stdout.
	//
	// If Kind is unknown, logger construction returns [ErrNotFound].
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty" validate:"omitempty,oneof=otlp json text tint"`

	// Protocol selects the OTLP transport protocol.
	//
	// Supported values are "http" and "grpc". When empty, "http" is used.
	// This field only applies when Kind is "otlp".
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty" toml:"protocol,omitempty" validate:"omitempty,oneof=http grpc"`

	// URL is the destination endpoint for the selected Kind, when applicable.
	//
	// For "otlp" over HTTP, this is the required OTLP/HTTP logs endpoint URL.
	// For "otlp" over gRPC, this is the required collector host:port endpoint.
	// Standard OpenTelemetry endpoint environment variables are not used as fallbacks;
	// configure this value explicitly through go-service config.
	URL string `yaml:"url,omitempty" json:"url,omitempty" toml:"url,omitempty"`

	// Level is the minimum log level to emit.
	//
	// Values are interpreted as slog levels (for example "debug", "info", "warn", "error").
	// When unset, the effective default is info.
	//
	// Unknown values are rejected by NewLogger with ErrInvalidLevel.
	Level string `yaml:"level,omitempty" json:"level,omitempty" toml:"level,omitempty"`

	// BatchTimeout is the maximum delay between batched log record exports.
	//
	// A zero value keeps the OpenTelemetry SDK default. Negative values are
	// invalid. This field only applies when Kind is "otlp".
	BatchTimeout time.Duration `yaml:"batch_timeout,omitempty" json:"batch_timeout,omitempty" toml:"batch_timeout,omitempty" validate:"gte=0"`

	// ExportTimeout is the maximum duration allowed for a single batch export.
	//
	// A zero value keeps the OpenTelemetry SDK default. Negative values are
	// invalid. This field only applies when Kind is "otlp".
	ExportTimeout time.Duration `yaml:"export_timeout,omitempty" json:"export_timeout,omitempty" toml:"export_timeout,omitempty" validate:"gte=0"`

	// MaxQueueSize is the maximum number of log records buffered before older
	// records are dropped.
	//
	// A zero value keeps the OpenTelemetry SDK default, which is 2048.
	// Negative values are invalid. This field only applies when Kind is "otlp".
	MaxQueueSize int `yaml:"max_queue_size,omitempty" json:"max_queue_size,omitempty" toml:"max_queue_size,omitempty" validate:"gte=0"`

	// MaxExportBatchSize is the maximum number of log records exported in a single batch.
	//
	// A zero value keeps the OpenTelemetry SDK default, which is 512. Negative
	// values are invalid. This field only applies when Kind is "otlp".
	MaxExportBatchSize int `yaml:"max_export_batch_size,omitempty" json:"max_export_batch_size,omitempty" toml:"max_export_batch_size,omitempty" validate:"gte=0"`
}

// IsEnabled reports whether logging configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetKind returns the configured logger kind.
//
// A nil receiver returns an empty kind, which callers treat as no configured
// stdout format.
func (c *Config) GetKind() string {
	if c == nil {
		return ""
	}

	return c.Kind
}

// GetProtocol returns the configured OTLP transport protocol.
//
// A nil receiver or an empty value falls back to OTLP/HTTP.
func (c *Config) GetProtocol() string {
	if c == nil || strings.IsEmpty(c.Protocol) {
		return otlp.ProtocolHTTP
	}

	return c.Protocol
}
