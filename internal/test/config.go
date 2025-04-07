package test

import (
	"time"

	cache "github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	"github.com/alexfalkowski/go-service/crypto/ssh"
	"github.com/alexfalkowski/go-service/crypto/tls"
	sql "github.com/alexfalkowski/go-service/database/sql/config"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/header"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/token"
	ts "github.com/alexfalkowski/go-service/token/ssh"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

const timeout = 2 * time.Second

// NewToken for test.
func NewToken(kind, secret string) *token.Config {
	return &token.Config{
		Kind:       kind,
		Secret:     Path(secret),
		Subject:    "sub",
		Audience:   "aud",
		Issuer:     "iss",
		Expiration: "1h",
		KeyID:      "1234567890",
		SSH: &ts.Config{
			Key: &ts.Key{
				Name:   "test",
				Config: NewSSH(),
			},
			Keys: ts.Keys{
				&ts.Key{
					Name:   "test",
					Config: NewSSH(),
				},
			},
		},
	}
}

// NewEd25519 for test.
func NewEd25519() *ed25519.Config {
	return &ed25519.Config{
		Public:  Path("secrets/ed25519_public"),
		Private: Path("secrets/ed25519_private"),
	}
}

// NewRSA for test.
func NewRSA() *rsa.Config {
	return &rsa.Config{
		Public:  Path("secrets/rsa_public"),
		Private: Path("secrets/rsa_private"),
	}
}

// NewAES for test.
func NewAES() *aes.Config {
	return &aes.Config{
		Key: Path("secrets/aes"),
	}
}

// NewHMAC for test.
func NewHMAC() *hmac.Config {
	return &hmac.Config{
		Key: Path("secrets/hmac"),
	}
}

// NewHook for test.
func NewHook() *hooks.Config {
	return &hooks.Config{
		Secret: Path("secrets/hooks"),
	}
}

// NewSSH for test.
func NewSSH() *ssh.Config {
	return &ssh.Config{
		Public:  Path("secrets/ssh_public"),
		Private: Path("secrets/ssh_private"),
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

// NewSecureClientConfig for test.
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
		Cert: Path(c),
		Key:  Path(k),
	}

	return tc
}

// NewInsecureTransportConfig for test.
func NewInsecureTransportConfig() *transport.Config {
	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Timeout: timeout.String(),
				Address: ":11000",
				Retry:   NewRetry(),
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout: timeout.String(),
				Address: ":12000",
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
				Address: ":11443",
				Retry:   retry,
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout: timeout.String(),
				TLS:     config,
				Address: ":12443",
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
			"Authorization": Path("secrets/telemetry"),
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
			"Authorization": Path("secrets/telemetry"),
		},
	}
}

// NewOTLPTracerConfig for test.
func NewOTLPTracerConfig() *tracer.Config {
	return &tracer.Config{
		Kind: "otlp",
		URL:  "http://localhost:4318/v1/traces",
		Headers: header.Map{
			"Authorization": Path("secrets/telemetry"),
		},
	}
}

// NewPGConfig for test.
func NewPGConfig() *pg.Config {
	return &pg.Config{
		Config: &sql.Config{
			Masters:         []sql.DSN{{URL: Path("secrets/pg")}},
			Slaves:          []sql.DSN{{URL: Path("secrets/pg")}},
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour.String(),
		},
	}
}

// NewInputConfig for test.
func NewInputConfig(set *cmd.FlagSet) *cmd.InputConfig {
	return cmd.NewInputConfig(set, Encoder, FS)
}

// NewOutputConfig for test.
func NewOutputConfig(set *cmd.FlagSet) *cmd.OutputConfig {
	return cmd.NewOutputConfig(set, Encoder, FS)
}

// NewInsecureDebugConfig for test.
func NewInsecureDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
			Timeout: "5s",
			Address: ":13000",
			Retry:   NewRetry(),
		},
	}
}

// NewInsecureDebugConfig for test.
func NewSecureDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
			Timeout: "5s",
			TLS:     NewTLSServerConfig(),
			Address: ":13443",
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
			"url": Path("secrets/" + secret),
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
