package test

import (
	"time"

	cache "github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	sql "github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/server"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
)

const timeout = 2 * time.Second

// Validator for testing.
var Validator = config.NewValidator()

// ConfigOptions for testing.
var ConfigOptions = map[string]string{
	"read_timeout":        "10m",
	"write_timeout":       "10m",
	"idle_timeout":        "10m",
	"read_header_timeout": "10m",
}

// NewAccessConfig for test.
func NewAccessConfig() *access.Config {
	return &access.Config{
		Policy: Path("configs/rbac.csv"),
	}
}

// NewDecoder for test.
func NewDecoder(set *flag.FlagSet) config.Decoder {
	decoder := config.NewDecoder(config.DecoderParams{
		Name:    Name,
		Flags:   set,
		Encoder: Encoder,
		FS:      FS,
	})

	return decoder
}

// NewToken for test.
func NewToken(kind string) *token.Config {
	return &token.Config{
		Kind: kind,
		JWT: &jwt.Config{
			Issuer:     "iss",
			Expiration: "1h",
			KeyID:      "1234567890",
		},
		Paseto: &paseto.Config{
			Issuer:     "iss",
			Expiration: "1h",
		},
		SSH: &ssh.Config{
			Key: &ssh.Key{
				Name:   UserID.String(),
				Config: NewSSH("secrets/ssh_public", "secrets/ssh_private"),
			},
			Keys: ssh.Keys{
				&ssh.Key{
					Name:   UserID.String(),
					Config: NewSSH("secrets/ssh_public", "secrets/ssh_private"),
				},
			},
		},
	}
}

// NewEd25519 for test.
func NewEd25519() *ed25519.Config {
	return &ed25519.Config{
		Public:  FilePath("secrets/ed25519_public"),
		Private: FilePath("secrets/ed25519_private"),
	}
}

// NewRSA for test.
func NewRSA() *rsa.Config {
	return &rsa.Config{
		Public:  FilePath("secrets/rsa_public"),
		Private: FilePath("secrets/rsa_private"),
	}
}

// NewAES for test.
func NewAES() *aes.Config {
	return &aes.Config{
		Key: FilePath("secrets/aes"),
	}
}

// NewHMAC for test.
func NewHMAC() *hmac.Config {
	return &hmac.Config{
		Key: FilePath("secrets/hmac"),
	}
}

// NewHook for test.
func NewHook() *hooks.Config {
	return &hooks.Config{
		Secret: FilePath("secrets/hooks"),
	}
}

// NewRetry for test.
func NewRetry() *retry.Config {
	return &retry.Config{
		Timeout:  timeout.String(),
		Backoff:  "100ms",
		Attempts: 1,
	}
}

// NewTLSClientConfig for test.
func NewTLSClientConfig() *tls.Config {
	return NewTLSConfig("certs/client-cert.pem", "certs/client-key.pem")
}

// NewInsecureConfig for test.
func NewInsecureConfig() *tls.Config {
	return &tls.Config{}
}

// NewTLSServerConfig for test.
func NewTLSServerConfig() *tls.Config {
	return NewTLSConfig("certs/cert.pem", "certs/key.pem")
}

// NewTLSConfig for test.
func NewTLSConfig(c, k string) *tls.Config {
	tc := &tls.Config{
		Cert: FilePath(c),
		Key:  FilePath(k),
	}

	return tc
}

// NewInsecureTransportConfig for test.
func NewInsecureTransportConfig() *transport.Config {
	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Timeout: timeout.String(),
				Address: RandomAddress(),
				Retry:   NewRetry(),
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout: timeout.String(),
				Address: RandomAddress(),
				Retry:   NewRetry(),
			},
		},
	}
}

// NewSecureTransportConfig for test.
func NewSecureTransportConfig() *transport.Config {
	config := NewTLSServerConfig()
	retry := NewRetry()

	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Timeout: timeout.String(),
				TLS:     config,
				Address: RandomAddress(),
				Retry:   retry,
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout: timeout.String(),
				TLS:     config,
				Address: RandomAddress(),
				Retry:   retry,
			},
		},
	}
}

// NewOTLPLoggerConfig for test.
func NewOTLPLoggerConfig() *logger.Config {
	return &logger.Config{
		Kind:  "otlp",
		Level: "debug",
		URL:   "http://localhost:3100/loki/api/v1/push",
		Headers: header.Map{
			"Authorization": FilePath("secrets/telemetry"),
		},
	}
}

// NewTextLoggerConfig for test.
func NewTextLoggerConfig() *logger.Config {
	return &logger.Config{
		Kind:  "text",
		Level: "debug",
	}
}

// NewJSONLoggerConfig for test.
func NewJSONLoggerConfig() *logger.Config {
	return &logger.Config{
		Kind:  "json",
		Level: "debug",
	}
}

// NewTintLoggerConfig for test.
func NewTintLoggerConfig() *logger.Config {
	return &logger.Config{
		Kind:  "tint",
		Level: "debug",
	}
}

// NewPrometheusMetricsConfig for test.
func NewPrometheusMetricsConfig() *metrics.Config {
	return &metrics.Config{
		Kind: "prometheus",
	}
}

// NewOTLPMetricsConfig for test.
func NewOTLPMetricsConfig() *metrics.Config {
	return &metrics.Config{
		Kind: "otlp",
		URL:  "http://localhost:9009/otlp/v1/metrics",
		Headers: header.Map{
			"Authorization": FilePath("secrets/telemetry"),
		},
	}
}

// NewOTLPTracerConfig for test.
func NewOTLPTracerConfig() *tracer.Config {
	return &tracer.Config{
		Kind: "otlp",
		URL:  "http://localhost:4318/v1/traces",
		Headers: header.Map{
			"Authorization": FilePath("secrets/telemetry"),
		},
	}
}

// NewPGConfig for test.
func NewPGConfig() *pg.Config {
	return &pg.Config{
		Config: &sql.Config{
			Masters:         []sql.DSN{{URL: FilePath("secrets/pg")}},
			Slaves:          []sql.DSN{{URL: FilePath("secrets/pg")}},
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour.String(),
		},
	}
}

// NewInsecureDebugConfig for test.
func NewInsecureDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
			Timeout: "5s",
			Address: RandomAddress(),
			Retry:   NewRetry(),
		},
	}
}

// NewSecureDebugConfig for test.
func NewSecureDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
			Timeout: "5s",
			TLS:     NewTLSServerConfig(),
			Address: RandomAddress(),
			Retry:   NewRetry(),
		},
	}
}

// NewCacheConfig for test.
func NewCacheConfig(kind, compressor, encoder, secret string) *cache.Config {
	return &cache.Config{
		Kind:       kind,
		Compressor: compressor, Encoder: encoder,
		Options: map[string]any{
			"url": FilePath("secrets/" + secret),
		},
	}
}

// NewLimiterConfig for test.
func NewLimiterConfig(kind, interval string, tokens uint64) *limiter.Config {
	return &limiter.Config{
		Kind:     kind,
		Interval: interval,
		Tokens:   tokens,
	}
}
