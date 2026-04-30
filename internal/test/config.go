package test

import (
	cache "github.com/alexfalkowski/go-service/v2/cache/config"
	"github.com/alexfalkowski/go-service/v2/config"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	tls "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	sql "github.com/alexfalkowski/go-service/v2/database/sql/config"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
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

// Validator is the shared config validator used by test helpers.
var Validator = config.NewValidator()

// ConfigOptions contains long-lived server timeout defaults used in tests that
// decode server configs from option maps.
var ConfigOptions = options.Map{
	"read_timeout":        "10m",
	"write_timeout":       "10m",
	"idle_timeout":        "10m",
	"read_header_timeout": "10m",
}

// NewAccessConfig returns an access model and policy config backed by the RBAC fixtures in `test/configs`.
func NewAccessConfig() *access.Config {
	return &access.Config{
		Model:  FilePath("configs/rbac.conf"),
		Policy: FilePath("configs/rbac.csv"),
	}
}

// NewDecoder constructs a config decoder wired to the shared test name, encoder map, and filesystem fixtures.
func NewDecoder(set *flag.FlagSet) config.Decoder {
	decoder := config.NewDecoder(config.DecoderParams{
		Name:    Name,
		Flags:   set,
		Encoder: Encoder,
		FS:      FS,
	})

	return decoder
}

// NewToken returns a token config populated with the standard JWT, Paseto, and SSH test fixtures.
func NewToken(kind string) *token.Config {
	return &token.Config{
		Kind: kind,
		JWT: &jwt.Config{
			Issuer:     "iss",
			Expiration: time.Hour,
			KeyID:      "1234567890",
		},
		Paseto: &paseto.Config{
			Issuer:     "iss",
			Expiration: time.Hour,
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

// NewEd25519 returns Ed25519 key sources that point at the shared test secret fixtures.
func NewEd25519() *ed25519.Config {
	return &ed25519.Config{
		Public:  FilePath("secrets/ed25519_public"),
		Private: FilePath("secrets/ed25519_private"),
	}
}

// NewRSA returns RSA key sources that point at the shared test secret fixtures.
func NewRSA() *rsa.Config {
	return &rsa.Config{
		Public:  FilePath("secrets/rsa_public"),
		Private: FilePath("secrets/rsa_private"),
	}
}

// NewAES returns an AES config backed by the shared symmetric key fixture.
func NewAES() *aes.Config {
	return &aes.Config{
		Key: FilePath("secrets/aes"),
	}
}

// NewHMAC returns an HMAC config backed by the shared MAC key fixture.
func NewHMAC() *hmac.Config {
	return &hmac.Config{
		Key: FilePath("secrets/hmac"),
	}
}

// NewHook returns webhook verification config backed by the shared hook secret fixture.
func NewHook() *hooks.Config {
	return &hooks.Config{
		Secret: FilePath("secrets/hooks"),
	}
}

// NewRetry returns a short retry policy suitable for deterministic tests.
func NewRetry() *retry.Config {
	return &retry.Config{
		Timeout:  timeout,
		Backoff:  100 * time.Millisecond,
		Attempts: 1,
	}
}

// NewTLSClientConfig returns the client certificate fixture used by secure transport tests.
func NewTLSClientConfig() *tls.Config {
	return NewTLSConfig("certs/client-cert.pem", "certs/client-key.pem")
}

// NewInsecureConfig returns an empty TLS config, which the transport treats as insecure mode.
func NewInsecureConfig() *tls.Config {
	return &tls.Config{}
}

// NewTLSServerConfig returns the server certificate fixture used by secure transport tests.
func NewTLSServerConfig() *tls.Config {
	return NewTLSConfig("certs/cert.pem", "certs/key.pem")
}

// NewTLSConfig returns a TLS config that resolves the given certificate and key relative to `test/`.
func NewTLSConfig(c, k string) *tls.Config {
	tc := &tls.Config{
		Cert: FilePath(c),
		Key:  FilePath(k),
	}

	return tc
}

// NewInsecureTransportConfig returns HTTP and gRPC transport configs that listen on ephemeral local addresses without TLS.
func NewInsecureTransportConfig() *transport.Config {
	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Timeout: timeout,
				Address: RandomAddress(),
				Retry:   NewRetry(),
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout: timeout,
				Address: RandomAddress(),
				Retry:   NewRetry(),
			},
		},
	}
}

// NewHTTPTransportConfig returns a minimal HTTP transport config for tests that fill fields manually.
func NewHTTPTransportConfig() *http.Config {
	return &http.Config{Config: &server.Config{}}
}

// NewGRPCTransportConfig returns a minimal gRPC transport config for tests that fill fields manually.
func NewGRPCTransportConfig() *grpc.Config {
	return &grpc.Config{Config: &server.Config{}}
}

// NewSecureTransportConfig returns HTTP and gRPC transport configs wired with the shared TLS server fixtures on ephemeral local addresses.
func NewSecureTransportConfig() *transport.Config {
	config := NewTLSServerConfig()
	retry := NewRetry()

	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Timeout: timeout,
				TLS:     config,
				Address: RandomAddress(),
				Retry:   retry,
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout: timeout,
				TLS:     config,
				Address: RandomAddress(),
				Retry:   retry,
			},
		},
	}
}

// NewOTLPLoggerConfig returns the OTLP logger config used by telemetry-heavy integration tests.
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

// NewTextLoggerConfig returns a human-readable debug logger config.
func NewTextLoggerConfig() *logger.Config {
	return &logger.Config{
		Kind:  "text",
		Level: "debug",
	}
}

// NewJSONLoggerConfig returns a structured JSON debug logger config.
func NewJSONLoggerConfig() *logger.Config {
	return &logger.Config{
		Kind:  "json",
		Level: "debug",
	}
}

// NewTintLoggerConfig returns a tint logger config for local debugging-oriented tests.
func NewTintLoggerConfig() *logger.Config {
	return &logger.Config{
		Kind:  "tint",
		Level: "debug",
	}
}

// NewPrometheusMetricsConfig returns the Prometheus metrics config used by the default world setup.
func NewPrometheusMetricsConfig() *metrics.Config {
	return &metrics.Config{
		Kind: "prometheus",
	}
}

// NewOTLPMetricsConfig returns the OTLP metrics exporter config backed by the shared telemetry secret fixture.
func NewOTLPMetricsConfig() *metrics.Config {
	return &metrics.Config{
		Kind: "otlp",
		URL:  "http://localhost:9009/otlp/v1/metrics",
		Headers: header.Map{
			"Authorization": FilePath("secrets/telemetry"),
		},
	}
}

// NewOTLPTracerConfig returns the OTLP trace exporter config backed by the shared telemetry secret fixture.
func NewOTLPTracerConfig() *tracer.Config {
	return &tracer.Config{
		Kind: "otlp",
		URL:  "http://localhost:4318/v1/traces",
		Headers: header.Map{
			"Authorization": FilePath("secrets/telemetry"),
		},
	}
}

// NewPGConfig returns the Postgres config used by database-backed integration tests.
func NewPGConfig() *pg.Config {
	return &pg.Config{
		Config: &sql.Config{
			Masters:         []sql.DSN{{URL: FilePath("secrets/pg")}},
			Slaves:          []sql.DSN{{URL: FilePath("secrets/pg")}},
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour,
		},
	}
}

// NewInsecureDebugConfig returns a debug server config bound to an ephemeral local address without TLS.
func NewInsecureDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
			Timeout: 5 * time.Second,
			Address: RandomAddress(),
			Retry:   NewRetry(),
		},
	}
}

// NewSecureDebugConfig returns a debug server config bound to an ephemeral local address with the shared TLS fixture.
func NewSecureDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
			Timeout: 5 * time.Second,
			TLS:     NewTLSServerConfig(),
			Address: RandomAddress(),
			Retry:   NewRetry(),
		},
	}
}

// NewCacheConfig returns a cache driver config that resolves its backend URL from a fixture secret under `test/secrets`.
func NewCacheConfig(kind, compressor, encoder, secret string) *cache.Config {
	return &cache.Config{
		Kind:       kind,
		Compressor: compressor, Encoder: encoder,
		Options: map[string]any{
			"url": FilePath("secrets/" + secret),
		},
	}
}

// NewLimiterConfig returns a limiter config with the supplied backend kind, refill interval, and token count.
func NewLimiterConfig(kind, interval string, tokens uint64) *limiter.Config {
	return &limiter.Config{
		Kind:     kind,
		Interval: time.MustParseDuration(interval),
		Tokens:   tokens,
	}
}
