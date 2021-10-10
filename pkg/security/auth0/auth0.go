package auth0

import (
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/security/jwt"
	"github.com/dgraph-io/ristretto"
)

// NewGenerator for Auth0.
func NewGenerator(cfg *Config, client *http.Client, cache *ristretto.Cache) jwt.Generator {
	var generator jwt.Generator = &generator{cfg: cfg, client: client}

	generator = &cachedGenerator{cfg: cfg, cache: cache, Generator: generator}

	return generator
}

// NewCertificator for Auth0.
func NewCertificator(cfg *Config, client *http.Client, cache *ristretto.Cache) Certificator {
	var certificator Certificator = &pem{cfg: cfg, client: client}

	certificator = &cachedPEM{cfg: cfg, cache: cache, Certificator: certificator}

	return certificator
}

// NewVerifier for Auth0.
func NewVerifier(cfg *Config, cert Certificator) jwt.Verifier {
	var verifier jwt.Verifier = &verifier{cfg: cfg, cert: cert}

	return verifier
}
