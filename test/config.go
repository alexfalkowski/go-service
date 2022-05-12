package test

import (
	"time"

	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport/grpc"
	gretry "github.com/alexfalkowski/go-service/transport/grpc/retry"
	"github.com/alexfalkowski/go-service/transport/http"
	hretry "github.com/alexfalkowski/go-service/transport/http/retry"
	"github.com/alexfalkowski/go-service/transport/nsq"
	nretry "github.com/alexfalkowski/go-service/transport/nsq/retry"
)

const timeout = 2 * time.Second

// NewGRPCConfig for test.
func NewGRPCConfig() *grpc.Config {
	return &grpc.Config{
		Port:      GenerateRandomPort(),
		UserAgent: "TestGRPC/1.0",
		Retry: gretry.Config{
			Timeout:  timeout,
			Attempts: 1,
		},
	}
}

// NewHTTPConfig for test.
func NewHTTPConfig() *http.Config {
	return &http.Config{
		Port:      GenerateRandomPort(),
		UserAgent: "TestHTTP/1.0",
		Retry: hretry.Config{
			Timeout:  timeout,
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
			Timeout:  timeout,
			Attempts: 1,
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
