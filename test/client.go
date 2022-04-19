package test

import (
	"net/http"

	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// NewHTTPClient for test.
func NewHTTPClient(logger *zap.Logger, tracer opentracing.Tracer) *http.Client {
	return NewHTTPClientWithRoundTripper(logger, tracer, nil)
}

// NewHTTPClientWithRoundTripper for test.
func NewHTTPClientWithRoundTripper(logger *zap.Logger, tracer opentracing.Tracer, roundTripper http.RoundTripper) *http.Client {
	return shttp.NewClient(
		shttp.WithClientConfig(NewHTTPConfig()), shttp.WithClientLogger(logger),
		shttp.WithClientRoundTripper(roundTripper), shttp.WithClientBreaker(),
		shttp.WithClientTracer(tracer), shttp.WithClientRetry(),
	)
}
