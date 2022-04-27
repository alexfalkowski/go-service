package test

import (
	"time"

	"github.com/alexfalkowski/go-service/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	"github.com/alexfalkowski/go-service/transport/grpc"
	gretry "github.com/alexfalkowski/go-service/transport/grpc/retry"
	"github.com/alexfalkowski/go-service/transport/http"
	hretry "github.com/alexfalkowski/go-service/transport/http/retry"
	"github.com/alexfalkowski/go-service/transport/nsq"
	nretry "github.com/alexfalkowski/go-service/transport/nsq/retry"
)

// NewGRPCConfig for test.
func NewGRPCConfig() *grpc.Config {
	return &grpc.Config{
		Port:      GenerateRandomPort(),
		UserAgent: "TestGRPC/1.0",
		Retry: gretry.Config{
			Timeout:  2 * time.Second, // nolint:gomnd
			Attempts: 1,
		},
	}
}

// NewHTTPConfig for test.
func NewHTTPConfig() *http.Config {
	return &http.Config{
		UserAgent: "TestHTTP/1.0",
		Retry: hretry.Config{
			Timeout:  2 * time.Second, // nolint:gomnd
			Attempts: 1,
		},
	}
}

// NewNSQConfig for test.
func NewNSQConfig() *nsq.Config {
	return &nsq.Config{
		LookupHost: "localhost:4161",
		Host:       "localhost:4150",
		UserAgent:  "TestNSQ/1.0",
		Retry: nretry.Config{
			Timeout:  2 * time.Second, // nolint:gomnd
			Attempts: 1,
		},
	}
}

// NewJaegerConfig for test.
func NewJaegerConfig() *jaeger.Config {
	return &jaeger.Config{
		Host: "localhost:6831",
	}
}

// NewDatadogConfig for test.
func NewDatadogConfig() *datadog.Config {
	return &datadog.Config{
		Host: "localhost:8126",
	}
}
