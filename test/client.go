package test

import (
	"net/http"

	shttp "github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/http/metrics/prometheus"
	"github.com/alexfalkowski/go-service/version"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// NewHTTPClient for test.
func NewHTTPClient(logger *zap.Logger, tracer opentracing.Tracer, version version.Version, metrics *prometheus.ClientMetrics) *http.Client {
	return NewHTTPClientWithRoundTripper(logger, tracer, version, metrics, nil)
}

// NewHTTPClientWithRoundTripper for test.
func NewHTTPClientWithRoundTripper(
	logger *zap.Logger, tracer opentracing.Tracer,
	version version.Version, metrics *prometheus.ClientMetrics,
	roundTripper http.RoundTripper,
) *http.Client {
	return shttp.NewClient(
		shttp.ClientParams{Version: version, Config: NewHTTPConfig()},
		shttp.WithClientLogger(logger),
		shttp.WithClientRoundTripper(roundTripper), shttp.WithClientBreaker(),
		shttp.WithClientTracer(tracer), shttp.WithClientRetry(),
		shttp.WithClientMetrics(metrics),
	)
}
