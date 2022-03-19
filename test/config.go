package test

import (
	"time"

	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/grpc/ratelimit"
	grpcRetry "github.com/alexfalkowski/go-service/transport/grpc/retry"
	"github.com/alexfalkowski/go-service/transport/http"
	httpRetry "github.com/alexfalkowski/go-service/transport/http/retry"
	"github.com/alexfalkowski/go-service/transport/nsq"
)

// NewGRPCConfig for test.
func NewGRPCConfig() *grpc.Config {
	return &grpc.Config{
		Port:      GenerateRandomPort(),
		UserAgent: "TestGRPC/1.0",
		Retry: grpcRetry.Config{
			Timeout:  2 * time.Second, // nolint:gomnd
			Attempts: 1,
		},
		RateLimit: ratelimit.Config{
			Every: 1 * time.Minute,
			Burst: 1,
		},
	}
}

// NewHTTPConfig for test.
func NewHTTPConfig() *http.Config {
	return &http.Config{
		UserAgent: "TestHTTP/1.0",
		Retry: httpRetry.Config{
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
	}
}
