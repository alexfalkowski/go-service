package test

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/database/sql/config"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/retry"
	sr "github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

const timeout = 2 * time.Second

// NewHook for test.
func NewHook() *hooks.Config {
	return &hooks.Config{
		Secret: "YWJjZGUxMjM0NQ==",
	}
}

// NewRetry for test.
func NewRetry() *retry.Config {
	return &retry.Config{
		Timeout:  timeout.String(),
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
	_, b, _, _ := runtime.Caller(0) //nolint:dogsled
	dir := filepath.Dir(b)

	cert, err := os.ReadFile(filepath.Join(dir, c))
	sr.Must(err)

	key, err := os.ReadFile(filepath.Join(dir, k))
	sr.Must(err)

	tc := &tls.Config{
		Cert: base64.StdEncoding.EncodeToString(cert),
		Key:  base64.StdEncoding.EncodeToString(key),
	}

	return tc
}

// NewInsecureTransportConfig for test.
func NewInsecureTransportConfig() *transport.Config {
	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Port:      Port(),
				UserAgent: "TestHTTP/1.0",
				Retry:     NewRetry(),
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
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
				TLS:       s,
				Port:      Port(),
				UserAgent: "TestHTTP/1.0",
				Retry:     r,
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
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
	}
}

// NewOTLPTracerConfig for test.
func NewOTLPTracerConfig() *tracer.Config {
	return &tracer.Config{
		Kind: "otlp",
		Host: "localhost:4318",
	}
}

// NewBaselimeTracerConfig for test.
func NewBaselimeTracerConfig() *tracer.Config {
	return &tracer.Config{
		Kind: "baselime",
		Key:  os.Getenv("BASELIME_API_KEY"),
	}
}

// NewPGConfig for test.
func NewPGConfig() *pg.Config {
	return &pg.Config{Config: &config.Config{
		Masters:         []config.DSN{{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"}},
		Slaves:          []config.DSN{{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"}},
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour.String(),
	}}
}

// NewCmdConfig for test.
func NewCmdConfig(flag string) (*cmd.Config, error) {
	return cmd.NewConfig(flag, Marshaller)
}

// NewInsecureDebugConfig for test.
func NewInsecureDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
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
			TLS:       NewTLSServerConfig(),
			Port:      Port(),
			UserAgent: "TestHTTPDebug/1.0",
			Retry:     NewRetry(),
		},
	}
}

// NewRedisConfig for test.
func NewRedisConfig(host, compressor, marshaller string) *redis.Config {
	return &redis.Config{
		Addresses:  map[string]string{"server": host},
		Compressor: compressor, Marshaller: marshaller,
	}
}

// NewLimiterConfig for test.
func NewLimiterConfig(pattern string) *limiter.Config {
	return &limiter.Config{
		Kind:    "user-agent",
		Pattern: pattern,
	}
}
