package test

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alexfalkowski/go-service/cmd"
	"github.com/alexfalkowski/go-service/database/sql/config"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/events"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

const timeout = 2 * time.Second

// Config for test.
type Config struct {
	Events        events.Config `yaml:"events,omitempty" json:"events,omitempty" toml:"events,omitempty"`
	config.Config `yaml:",inline" json:",inline" toml:",inline"`
}

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

// NewSecureClientConfig for test.
func NewSecureClientConfig() *security.Config {
	_, b, _, _ := runtime.Caller(0) //nolint:dogsled
	dir := filepath.Dir(b)

	return &security.Config{
		Enabled:  true,
		CertFile: filepath.Join(dir, "certs/client-cert.pem"),
		KeyFile:  filepath.Join(dir, "certs/client-key.pem"),
	}
}

// NewInsecureTransportConfig for test.
func NewInsecureTransportConfig() *transport.Config {
	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Enabled:   true,
				Port:      Port(),
				UserAgent: "TestHTTP/1.0",
				Retry:     NewRetry(),
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Enabled:   true,
				Port:      Port(),
				UserAgent: "TestGRPC/1.0",
				Retry:     NewRetry(),
			},
		},
	}
}

// NewSecureTransportConfig for test.
func NewSecureTransportConfig() *transport.Config {
	_, b, _, _ := runtime.Caller(0) //nolint:dogsled
	dir := filepath.Dir(b)

	s := &security.Config{
		Enabled:  true,
		CertFile: filepath.Join(dir, "certs/cert.pem"),
		KeyFile:  filepath.Join(dir, "certs/key.pem"),
	}
	r := NewRetry()

	return &transport.Config{
		HTTP: &http.Config{
			Config: &server.Config{
				Enabled:   true,
				Security:  s,
				Port:      Port(),
				UserAgent: "TestHTTP/1.0",
				Retry:     r,
			},
		},
		GRPC: &grpc.Config{
			Config: &server.Config{
				Enabled:   true,
				Security:  s,
				Port:      Port(),
				UserAgent: "TestGRPC/1.0",
				Retry:     r,
			},
		},
	}
}

// NewDefaultTracerConfig for test.
func NewDefaultTracerConfig() *tracer.Config {
	return &tracer.Config{
		Enabled: true,
		Host:    "localhost:4318",
	}
}

// NewBaselimeTracerConfig for test.
func NewBaselimeTracerConfig() *tracer.Config {
	return &tracer.Config{
		Enabled: true,
		Kind:    "baselime",
		Key:     os.Getenv("BASELIME_API_KEY"),
	}
}

// NewPGConfig for test.
func NewPGConfig() *pg.Config {
	return &pg.Config{Config: config.Config{
		Masters:         []config.DSN{{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"}},
		Slaves:          []config.DSN{{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"}},
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour.String(),
	}}
}

// NewCmdConfig for test.
func NewCmdConfig(flag string) (*cmd.Config, error) {
	p := marshaller.FactoryParams{YAML: marshaller.NewYAML(), TOML: marshaller.NewTOML()}

	return cmd.NewConfig(flag, marshaller.NewFactory(p))
}

// NewDebugConfig for test.
func NewDebugConfig() *debug.Config {
	return &debug.Config{
		Config: &server.Config{
			Enabled:   true,
			Port:      Port(),
			UserAgent: "TestHTTPDebug/1.0",
			Retry:     NewRetry(),
		},
	}
}
