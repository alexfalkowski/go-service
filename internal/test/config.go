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
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/telemetry/header"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/keys"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	"github.com/alexfalkowski/go-service/v2/transport"
	"github.com/alexfalkowski/go-service/v2/transport/breaker"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	grpcretry "github.com/alexfalkowski/go-service/v2/transport/grpc/retry"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	httpretry "github.com/alexfalkowski/go-service/v2/transport/http/retry"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/retry"
)

const timeout = 2 * time.Second

// Validator is the shared config validator used by test helpers.
var Validator = config.NewValidator()

// FastRetryConfig is a shared client retry config for tests that need one retry with minimal backoff.
var FastRetryConfig = &retry.Config{
	Backoff:  time.Nanosecond,
	Attempts: 2,
}

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
			Key:        "1234567890",
			Keys: keys.Map{
				"1234567890": {
					Config: NewEd25519(),
				},
			},
		},
		Paseto: &paseto.Config{
			Issuer:     "iss",
			Expiration: time.Hour,
			Key:        "1234567890",
			Keys: keys.Map{
				"1234567890": {
					Config: NewEd25519(),
				},
			},
		},
		SSH: &ssh.Config{
			Expiration: time.Hour,
			Key:        UserID.String(),
			Keys: ssh.Keys{
				UserID.String(): {
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
		Key: "current",
		Secrets: hooks.Secrets{
			"current": FilePath("secrets/hooks"),
		},
	}
}

// NewRetry returns the default client retry config used by the test world.
//
// Attempts is set to 1, so callers exercise retry middleware wiring without
// scheduling additional attempts by default.
func NewRetry() *retry.Config {
	return &retry.Config{
		Timeout:  timeout,
		Backoff:  100 * time.Millisecond,
		Attempts: 1,
	}
}

// NewBreaker returns a client breaker config with a consecutive failure threshold.
func NewBreaker(consecutiveFailures uint32) *breaker.Config {
	return &breaker.Config{ConsecutiveFailures: consecutiveFailures}
}

// NewHTTPRetryConfig returns an HTTP retry config with shared retry mechanics.
func NewHTTPRetryConfig(attempts uint64, backoff time.Duration, statusCodes ...int) *httpretry.Config {
	return httpretry.NewConfig(&retry.Config{Attempts: attempts, Backoff: backoff}, statusCodes...)
}

// NewGRPCRetryConfig returns a gRPC retry config with shared retry mechanics and a per-attempt timeout.
func NewGRPCRetryConfig(attempts uint64, backoff time.Duration, codes ...codes.Code) *grpcretry.Config {
	return grpcretry.NewConfig(&retry.Config{Attempts: attempts, Timeout: time.Second, Backoff: backoff}, codes...)
}

// NewTLSClientConfig returns the client certificate fixture used by secure transport tests.
func NewTLSClientConfig() *tls.Config {
	return NewTLSConfig("certs/client-cert.pem", "certs/client-key.pem")
}

// NewInsecureConfig returns an empty TLS config for client tests that enable TLS without loading key material.
func NewInsecureConfig() *tls.Config {
	return &tls.Config{}
}

// NewTLSServerConfig returns the server certificate fixture used by secure transport tests.
func NewTLSServerConfig() *tls.Config {
	return NewTLSConfig("certs/cert.pem", "certs/key.pem")
}

// NewTLSConfig returns a TLS config that resolves test TLS fixtures relative to `test/`.
func NewTLSConfig(c, k string) *tls.Config {
	tc := &tls.Config{
		Cert:       FilePath(c),
		Key:        FilePath(k),
		CA:         FilePath("certs/rootCA.pem"),
		ServerName: "localhost",
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
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout: timeout,
				Address: RandomAddress(),
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

	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Timeout: timeout,
				TLS:     config,
				Address: RandomAddress(),
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout: timeout,
				TLS:     config,
				Address: RandomAddress(),
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
	return NewPGConfigWithDSNs(
		[]sql.DSN{{URL: FilePath("secrets/pg")}},
		[]sql.DSN{{URL: FilePath("secrets/pg")}},
	)
}

// NewPGConfigWithDSNs returns the Postgres config used by database-backed integration tests with explicit DSNs.
func NewPGConfigWithDSNs(readers, writers []sql.DSN) *pg.Config {
	return &pg.Config{
		Config: &sql.Config{
			Reader: newSQLPool(readers),
			Writer: newSQLPool(writers),
		},
	}
}

func newSQLPool(dsns []sql.DSN) *sql.Pool {
	return &sql.Pool{
		DSNs: dsns,
		Settings: &sql.PoolSettings{
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxIdleTime: 30 * time.Minute,
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
		},
	}
}

// NewCacheConfig returns a cache driver config that resolves its backend URL from a fixture secret under `test/secrets`.
func NewCacheConfig(kind, compressor, encoder, secret string) *cache.Config {
	return &cache.Config{
		Kind:       kind,
		Compressor: compressor, Encoder: encoder,
		MaxEntries: cache.DefaultMaxEntries,
		Options: map[string]any{
			"url": FilePath("secrets/" + secret),
		},
	}
}

// NewLimiterConfig returns a limiter config with the supplied key kind, refill interval, and token count.
func NewLimiterConfig(kind, interval string, tokens uint64) *limiter.Config {
	return &limiter.Config{
		Kind:     kind,
		Interval: time.MustParseDuration(interval),
		Tokens:   tokens,
	}
}
