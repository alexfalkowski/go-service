package test

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/database/sql/config"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

const timeout = 2 * time.Second

// NewEd25519 for test.
func NewEd25519() *ed25519.Config {
	return &ed25519.Config{
		Public:  ed25519.PublicKey(hooks.Secret(Path("secrets/ed25519_public"))),
		Private: ed25519.PrivateKey(hooks.Secret(Path("secrets/ed25519_private"))),
	}
}

// NewEd25519 for test.
func NewRSA() *rsa.Config {
	return &rsa.Config{
		Public:  rsa.PublicKey(hooks.Secret(Path("secrets/rsa_public"))),
		Private: rsa.PrivateKey(hooks.Secret(Path("secrets/rsa_private"))),
	}
}

// NewAES for test.
func NewAES() *aes.Config {
	return &aes.Config{
		Key: aes.Key(Path("secrets/aes")),
	}
}

// NewHMAC for test.
func NewHMAC() *hmac.Config {
	return &hmac.Config{
		Key: hmac.Key(Path("secrets/hmac")),
	}
}

// NewHook for test.
func NewHook() *hooks.Config {
	return &hooks.Config{
		Secret: hooks.Secret(Path("secrets/hooks")),
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
		Cert: tls.Cert(Path(c)),
		Key:  tls.Key(Path(k)),
	}

	return tc
}

// NewInsecureTransportConfig for test.
func NewInsecureTransportConfig() *transport.Config {
	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Timeout:   timeout.String(),
				Port:      Port(),
				UserAgent: "TestHTTP/1.0",
				Retry:     NewRetry(),
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout:   timeout.String(),
				Port:      Port(),
				UserAgent: "TestGRPC/1.0",
				Retry:     NewRetry(),
			},
		},
	}
}

// NewSecureTransportConfig for test.
func NewSecureTransportConfig() *transport.Config {
	s := NewTLSServerConfig()
	r := NewRetry()

	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Timeout:   timeout.String(),
				TLS:       s,
				Port:      Port(),
				UserAgent: "TestHTTP/1.0",
				Retry:     r,
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout:   timeout.String(),
				TLS:       s,
				Port:      Port(),
				UserAgent: "TestGRPC/1.0",
				Retry:     r,
			},
		},
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
		Host: "http://localhost:9009/otlp/v1/metrics",
		Key:  metrics.Key(Path("secrets/basic")),
	}
}

// NewOTLPTracerConfig for test.
func NewOTLPTracerConfig() *tracer.Config {
	return &tracer.Config{
		Kind: "otlp",
		Host: "localhost:4318",
		Key:  tracer.Key(Path("secrets/basic")),
	}
}

// NewPGConfig for test.
func NewPGConfig() *pg.Config {
	return &pg.Config{
		Config: &config.Config{
			Masters:         []config.DSN{{URL: config.URL(Path("secrets/pg"))}},
			Slaves:          []config.DSN{{URL: config.URL(Path("secrets/pg"))}},
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour.String(),
		},
	}
}

// NewInputConfig for test.
func NewInputConfig(flag string) *cmd.InputConfig {
	*cmd.InputFlag = flag

	return cmd.NewInputConfig(Marshaller)
}

// NewOutputConfig for test.
func NewOutputConfig(flag string) *cmd.OutputConfig {
	*cmd.OutputFlag = flag

	return cmd.NewOutputConfig(Marshaller)
}

// NewInsecureDebugConfig for test.
func NewInsecureDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
			Timeout:   "5s",
			Port:      Port(),
			UserAgent: "TestHTTPDebug/1.0",
			Retry:     NewRetry(),
		},
	}
}

// NewInsecureDebugConfig for test.
func NewSecureDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
			Timeout:   "5s",
			TLS:       NewTLSServerConfig(),
			Port:      Port(),
			UserAgent: "TestHTTPDebug/1.0",
			Retry:     NewRetry(),
		},
	}
}

// NewRedisConfig for test.
func NewRedisConfig(secret, compressor, marshaller string) *redis.Config {
	return &redis.Config{
		Addresses:  map[string]string{"server": "localhost:6379"},
		Compressor: compressor, Marshaller: marshaller,
		URL: redis.URL(Path("secrets/" + secret)),
	}
}

// NewLimiterConfig for test.
func NewLimiterConfig(kind, pattern string) *limiter.Config {
	return &limiter.Config{
		Kind:    kind,
		Pattern: pattern,
	}
}

// Path for test.
func Path(p string) string {
	_, b, _, _ := runtime.Caller(0) //nolint:dogsled
	dir := filepath.Dir(b)

	return filepath.Join(dir, p)
}
