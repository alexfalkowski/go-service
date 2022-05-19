package test

import (
	"time"

	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/cache/ristretto"
	"github.com/alexfalkowski/go-service/database/sql/config"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	"github.com/alexfalkowski/go-service/security/auth0"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport"
	"github.com/alexfalkowski/go-service/transport/grpc"
	gretry "github.com/alexfalkowski/go-service/transport/grpc/retry"
	"github.com/alexfalkowski/go-service/transport/http"
	hretry "github.com/alexfalkowski/go-service/transport/http/retry"
	"github.com/alexfalkowski/go-service/transport/nsq"
	nretry "github.com/alexfalkowski/go-service/transport/nsq/retry"
)

const timeout = 2 * time.Second

// Config for test.
type Config struct {
	Runtime string
}

func (cfg *Config) GetRuntime() string {
	return cfg.Runtime
}

func (cfg *Config) RedisConfig() *redis.Config {
	return nil
}

func (cfg *Config) RistrettoConfig() *ristretto.Config {
	return nil
}

func (cfg *Config) Auth0Config() *auth0.Config {
	return nil
}

func (cfg *Config) PGConfig() *pg.Config {
	return nil
}

func (cfg *Config) OpentracingConfig() *opentracing.Config {
	return nil
}

func (cfg *Config) TransportConfig() *transport.Config {
	return nil
}

func (cfg *Config) GRPCConfig() *grpc.Config {
	return nil
}

func (cfg *Config) HTTPConfig() *http.Config {
	return nil
}

func (cfg *Config) NSQConfig() *nsq.Config {
	return nil
}

// NewTransportConfig for test.
func NewTransportConfig() *transport.Config {
	return &transport.Config{
		Port: GenerateRandomPort(),
		HTTP: http.Config{
			UserAgent: "TestHTTP/1.0",
			Retry: hretry.Config{
				Timeout:  timeout,
				Attempts: 1,
			},
		},
		GRPC: grpc.Config{
			UserAgent: "TestGRPC/1.0",
			Retry: gretry.Config{
				Timeout:  timeout,
				Attempts: 1,
			},
		},
		NSQ: nsq.Config{
			LookupHost: "localhost:4161",
			Host:       "localhost:4150",
			UserAgent:  "TestNSQ/1.0",
			Retry: nretry.Config{
				Timeout:  timeout,
				Attempts: 1,
			},
		},
	}
}

// NewJaegerConfig for test.
func NewJaegerConfig() *opentracing.Config {
	return &opentracing.Config{
		Type: "jaeger",
		Host: "localhost:6831",
	}
}

// NewDatadogConfig for test.
func NewDatadogConfig() *opentracing.Config {
	return &opentracing.Config{
		Type: "datadog",
		Host: "localhost:8126",
	}
}

// NewPGConfig for test.
// nolint:gomnd
func NewPGConfig() *pg.Config {
	return &pg.Config{Config: config.Config{
		Masters:         []config.DSN{{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"}},
		Slaves:          []config.DSN{{URL: "postgres://test:test@localhost:5432/test?sslmode=disable"}},
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}}
}
