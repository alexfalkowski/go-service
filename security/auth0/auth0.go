package auth0

import (
	"github.com/alexfalkowski/go-service/security/jwt"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/dgraph-io/ristretto"
	"go.uber.org/zap"
)

// NewGenerator for Auth0.
// nolint:ireturn
func NewGenerator(cfg *Config, httpCfg *http.Config, logger *zap.Logger, cache *ristretto.Cache, tracer opentracing.TransportTracer) jwt.Generator {
	client := http.NewClient(
		http.WithClientConfig(httpCfg), http.WithClientLogger(logger),
		http.WithClientBreaker(), http.WithClientRetry(),
		http.WithClientTracer(tracer),
	)

	var generator jwt.Generator = &generator{cfg: cfg, client: client}
	generator = &cachedGenerator{cfg: cfg, cache: cache, Generator: generator}

	return generator
}

// NewCertificator for Auth0.
// nolint:ireturn
func NewCertificator(cfg *Config, httpCfg *http.Config, logger *zap.Logger, cache *ristretto.Cache, tracer opentracing.TransportTracer) Certificator {
	client := http.NewClient(
		http.WithClientConfig(httpCfg), http.WithClientLogger(logger),
		http.WithClientBreaker(), http.WithClientRetry(),
		http.WithClientTracer(tracer),
	)

	var certificator Certificator = &pem{cfg: cfg, client: client}
	certificator = &cachedPEM{cfg: cfg, cache: cache, Certificator: certificator}

	return certificator
}

// NewVerifier for Auth0.
// nolint:ireturn
func NewVerifier(cfg *Config, cert Certificator) jwt.Verifier {
	var verifier jwt.Verifier = &verifier{cfg: cfg, cert: cert}

	return verifier
}
