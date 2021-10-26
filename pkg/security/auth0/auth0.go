package auth0

import (
	"github.com/alexfalkowski/go-service/pkg/security/jwt"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/dgraph-io/ristretto"
	"go.uber.org/zap"
)

// NewGenerator for Auth0.
func NewGenerator(cfg *Config, logger *zap.Logger, cache *ristretto.Cache) jwt.Generator {
	var generator jwt.Generator = &generator{cfg: cfg, client: http.NewClient(&http.ClientParams{Logger: logger})}

	generator = &cachedGenerator{cfg: cfg, cache: cache, Generator: generator}

	return generator
}

// NewCertificator for Auth0.
func NewCertificator(cfg *Config, logger *zap.Logger, cache *ristretto.Cache) Certificator {
	var certificator Certificator = &pem{cfg: cfg, client: http.NewClient(&http.ClientParams{Logger: logger})}

	certificator = &cachedPEM{cfg: cfg, cache: cache, Certificator: certificator}

	return certificator
}

// NewVerifier for Auth0.
func NewVerifier(cfg *Config, cert Certificator) jwt.Verifier {
	var verifier jwt.Verifier = &verifier{cfg: cfg, cert: cert}

	return verifier
}
