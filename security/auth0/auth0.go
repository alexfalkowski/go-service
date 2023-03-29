package auth0

import (
	"github.com/alexfalkowski/go-service/security/jwt"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/alexfalkowski/go-service/transport/http/otel"
	"github.com/dgraph-io/ristretto"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// GeneratorParams for Auth0.
type GeneratorParams struct {
	fx.In

	Config     *Config
	HTTPConfig *http.Config
	Logger     *zap.Logger
	Cache      *ristretto.Cache
	Tracer     otel.Tracer
}

// NewGenerator for Auth0.
func NewGenerator(params GeneratorParams) jwt.Generator {
	client := http.NewClient(
		http.ClientParams{Config: params.HTTPConfig},
		http.WithClientLogger(params.Logger),
		http.WithClientBreaker(), http.WithClientRetry(),
		http.WithClientTracer(params.Tracer),
	)

	var generator jwt.Generator = &generator{cfg: params.Config, client: client}
	generator = &cachedGenerator{cfg: params.Config, cache: params.Cache, Generator: generator}

	return generator
}

// CertificatorParams for Auth0.
type CertificatorParams struct {
	fx.In

	Config     *Config
	HTTPConfig *http.Config
	Logger     *zap.Logger
	Cache      *ristretto.Cache
	Tracer     otel.Tracer
}

// NewCertificator for Auth0.
func NewCertificator(params CertificatorParams) Certificator {
	client := http.NewClient(
		http.ClientParams{Config: params.HTTPConfig},
		http.WithClientLogger(params.Logger),
		http.WithClientBreaker(), http.WithClientRetry(),
		http.WithClientTracer(params.Tracer),
	)

	var certificator Certificator = &pem{cfg: params.Config, client: client}
	certificator = &cachedPEM{cfg: params.Config, cache: params.Cache, Certificator: certificator}

	return certificator
}

// NewVerifier for Auth0.
func NewVerifier(cfg *Config, cert Certificator) jwt.Verifier {
	var verifier jwt.Verifier = &verifier{cfg: cfg, cert: cert}

	return verifier
}
