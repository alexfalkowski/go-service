package test

import (
	"net/http"

	"github.com/alexfalkowski/go-service/trace/opentracing"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	"go.uber.org/zap"
)

// NewHTTPClient for test.
func NewHTTPClient(logger *zap.Logger, tracer opentracing.TransportTracer) *http.Client {
	return NewHTTPClientWithRoundTripper(logger, tracer, nil)
}

// NewHTTPClientWithRoundTripper for test.
func NewHTTPClientWithRoundTripper(logger *zap.Logger, tracer opentracing.TransportTracer, roundTripper http.RoundTripper) *http.Client {
	return shttp.NewClient(
		shttp.WithClientConfig(NewHTTPConfig()), shttp.WithClientLogger(logger),
		shttp.WithClientRoundTripper(roundTripper), shttp.WithClientBreaker(),
		shttp.WithClientTracer(tracer), shttp.WithClientRetry(),
	)
}
