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
	"github.com/alexfalkowski/go-service/crypto/ssh"
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/database/sql/config"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/header"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

const timeout = 2 * time.Second

// NewToken for test.
func NewToken(kind string) *token.Config {
	return &token.Config{
		Kind:       kind,
		Key:        Path("secrets/token"),
		Subject:    "sub",
		Audience:   "aud",
		Issuer:     "iss",
		Expiration: "1h",
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
				Address: "localhost:" + Port(),
				Retry:   NewRetry(),
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout: timeout.String(),
				Address: "localhost:" + Port(),
				Retry:   NewRetry(),
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
				Timeout: timeout.String(),
				TLS:     s,
				Address: "localhost:" + Port(),
				Retry:   r,
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Timeout: timeout.String(),
				TLS:     s,
				Address: "localhost:" + Port(),
				Retry:   r,
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
		URL:  "http://localhost:9009/otlp/v1/metrics",
		Headers: header.Map{
			"Authorization": Path("secrets/metrics"),
		},
	}
}

// NewOTLPTracerConfig for test.
func NewOTLPTracerConfig() *tracer.Config {
	return &tracer.Config{
		Kind: "otlp",
		URL:  "localhost:4318",
		Headers: header.Map{
			"Authorization": Path("secrets/tracer"),
		},
	}
}

// NewPGConfig for test.
func NewPGConfig() *pg.Config {
	return &pg.Config{
		Config: &config.Config{
			Masters:         []config.DSN{{URL: Path("secrets/pg")}},
			Slaves:          []config.DSN{{URL: Path("secrets/pg")}},
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour.String(),
		},
	}
}

// NewInputConfig for test.
func NewInputConfig(flag string) *cmd.InputConfig {
	*cmd.InputFlag = flag

	return cmd.NewInputConfig(Encoder)
}

// NewOutputConfig for test.
func NewOutputConfig(flag string) *cmd.OutputConfig {
	*cmd.OutputFlag = flag

	return cmd.NewOutputConfig(Encoder)
}

// NewInsecureDebugConfig for test.
func NewInsecureDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
			Timeout: "5s",
			Address: "localhost:" + Port(),
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
			Address: "localhost:" + Port(),
			Retry:   NewRetry(),
		},
	}
}

// NewRedisConfig for test.
func NewRedisConfig(secret, compressor, encoder string) *redis.Config {
	return &redis.Config{
		Addresses:  map[string]string{"server": "localhost:6379"},
		Compressor: compressor, Encoder: encoder,
		URL: Path("secrets/" + secret),
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

// Path for test.
func Path(p string) string {
	_, b, _, _ := runtime.Caller(0) //nolint:dogsled
	dir := filepath.Dir(b)

	return filepath.Join(dir, p)
}
