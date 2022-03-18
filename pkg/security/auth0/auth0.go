package auth0

import (
	"github.com/alexfalkowski/go-service/pkg/security/jwt"
	"github.com/alexfalkowski/go-service/pkg/transport/http"
	"github.com/dgraph-io/ristretto"
	"go.uber.org/zap"
)

// NewGenerator for Auth0.
// nolint:ireturn
func NewGenerator(cfg *Config, httpCfg *http.Config, logger *zap.Logger, cache *ristretto.Cache) jwt.Generator {
	params := &http.ClientParams{
		Config: httpCfg,
		Logger: logger,
	}

	var generator jwt.Generator = &generator{cfg: cfg, client: http.NewClient(params)}

	generator = &cachedGenerator{cfg: cfg, cache: cache, Generator: generator}

	return generator
}

// NewCertificator for Auth0.
// nolint:ireturn
func NewCertificator(cfg *Config, httpCfg *http.Config, logger *zap.Logger, cache *ristretto.Cache) Certificator {
	params := &http.ClientParams{
		Config: httpCfg,
		Logger: logger,
	}

	var certificator Certificator = &pem{cfg: cfg, client: http.NewClient(params)}

	certificator = &cachedPEM{cfg: cfg, cache: cache, Certificator: certificator}

	return certificator
}

// NewVerifier for Auth0.
// nolint:ireturn
func NewVerifier(cfg *Config, cert Certificator) jwt.Verifier {
	var verifier jwt.Verifier = &verifier{cfg: cfg, cert: cert}

	return verifier
}
