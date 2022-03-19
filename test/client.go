package test

import (
	"net/http"

	shttp "github.com/alexfalkowski/go-service/transport/http"
	"go.uber.org/zap"
)

// NewHTTPClient for test.
func NewHTTPClient(logger *zap.Logger) *http.Client {
	return NewHTTPClientWithRoundTripper(logger, nil)
}

// NewHTTPClientWithRoundTripper for test.
func NewHTTPClientWithRoundTripper(logger *zap.Logger, roundTripper http.RoundTripper) *http.Client {
	params := &shttp.ClientParams{
		Config:       NewHTTPConfig(),
		Logger:       logger,
		RoundTripper: roundTripper,
	}

	return shttp.NewClient(params)
}
