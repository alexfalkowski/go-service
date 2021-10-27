package test

import (
	"net/http"

	pkgHTTP "github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/alexfalkowski/go-service/pkg/transport/http/retry"
	"go.uber.org/zap"
)

// NewHTTPConfig for test.
func NewHTTPConfig() *pkgHTTP.Config {
	return &pkgHTTP.Config{
		Retry: retry.Config{
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
