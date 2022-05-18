package test

import (
	"time"

	"github.com/alexfalkowski/go-service/database/sql/config"
	"github.com/alexfalkowski/go-service/database/sql/pg"
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
