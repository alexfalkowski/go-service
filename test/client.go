package test

import (
	"net/http"

	pkgGRPC "github.com/alexfalkowski/go-service/pkg/transport/grpc"
	pkgGRPCRetry "github.com/alexfalkowski/go-service/pkg/transport/grpc/retry"
	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	pkgHTTPRetry "github.com/alexfalkowski/go-service/pkg/transport/http/retry"
	"go.uber.org/zap"
)

// NewGRPCConfig for test.
func NewGRPCConfig() *pkgGRPC.Config {
	return &pkgGRPC.Config{
		Port: GenerateRandomPort(),
		Retry: pkgGRPCRetry.Config{
			Timeout:  2, // nolint:gomnd
			Attempts: 1,
		},
	}
}

// NewHTTPConfig for test.
func NewHTTPConfig() *pkgHTTP.Config {
	return &pkgHTTP.Config{
		Retry: pkgHTTPRetry.Config{
			Timeout:  2, // nolint:gomnd
			Attempts: 1,
		},
	}
}

// NewHTTPClient for test.
func NewHTTPClient(logger *zap.Logger) *http.Client {
	return NewHTTPClientWithRoundTripper(logger, nil)
}

// NewHTTPClientWithRoundTripper for test.
func NewHTTPClientWithRoundTripper(logger *zap.Logger, roundTripper http.RoundTripper) *http.Client {
	params := &pkgHTTP.ClientParams{
		Config:       NewHTTPConfig(),
		Logger:       logger,
		RoundTripper: roundTripper,
	}

	return pkgHTTP.NewClient(params)
}
