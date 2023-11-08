package oauth

import (
	"github.com/alexfalkowski/go-service/security/token"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"github.com/dgraph-io/ristretto"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// GeneratorParams for OAuth.
type GeneratorParams struct {
	fx.In

	Config     *Config
	HTTPConfig *http.Config
	Logger     *zap.Logger
	Cache      *ristretto.Cache
	Tracer     tracer.Tracer
	Meter      metric.Meter
}

// NewGenerator for OAuth.
func NewGenerator(params GeneratorParams) (token.Generator, error) {
	client, err := http.NewClient(params.HTTPConfig,
		http.WithClientLogger(params.Logger),
		http.WithClientBreaker(), http.WithClientRetry(),
		http.WithClientTracer(params.Tracer), http.WithClientMetrics(params.Meter),
	)
	if err != nil {
		return nil, err
	}

	var generator token.Generator = &generator{cfg: params.Config, client: client}
	generator = &cachedGenerator{cfg: params.Config, cache: params.Cache, Generator: generator}

	return generator, nil
}

// CertificatorParams for OAuth.
type CertificatorParams struct {
	fx.In

	Config     *Config
	HTTPConfig *http.Config
	Logger     *zap.Logger
	Cache      *ristretto.Cache
	Tracer     tracer.Tracer
	Meter      metric.Meter
}

// NewCertificator for OAuth.
func NewCertificator(params CertificatorParams) (Certificator, error) {
	client, err := http.NewClient(params.HTTPConfig,
		http.WithClientLogger(params.Logger),
		http.WithClientBreaker(), http.WithClientRetry(),
		http.WithClientTracer(params.Tracer), http.WithClientMetrics(params.Meter),
	)
	if err != nil {
		return nil, err
	}

	var certificator Certificator = &certificate{cfg: params.Config, client: client}
	certificator = &cachedCertificate{cfg: params.Config, cache: params.Cache, Certificator: certificator}

	return certificator, nil
}

// NewVerifier for OAuth.
func NewVerifier(cfg *Config, cert Certificator) token.Verifier {
	var verifier token.Verifier = &verifier{cfg: cfg, cert: cert}

	return verifier
}
