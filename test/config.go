package test

import (
	"github.com/alexfalkowski/go-service/pkg/transport/grpc"
	grpcRetry "github.com/alexfalkowski/go-service/pkg/transport/grpc/retry"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	httpRetry "github.com/alexfalkowski/go-service/pkg/transport/http/retry"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq"
)

// NewGRPCConfig for test.
func NewGRPCConfig() *grpc.Config {
	return &grpc.Config{
		Port:      GenerateRandomPort(),
		UserAgent: UserAgent,
		Retry: grpcRetry.Config{
			Timeout:  2, // nolint:gomnd
			Attempts: 1,
		},
	}
}

// NewHTTPConfig for test.
func NewHTTPConfig() *http.Config {
	return &http.Config{
		UserAgent: UserAgent,
		Retry: httpRetry.Config{
			Timeout:  2, // nolint:gomnd
			Attempts: 1,
		},
	}
}

// NewNSQConfig for test.
func NewNSQConfig() *nsq.Config {
	return &nsq.Config{
		LookupHost: "localhost:4161",
		Host:       "localhost:4150",
		UserAgent:  UserAgent,
	}
}
