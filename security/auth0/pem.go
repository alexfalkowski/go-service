package auth0

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/form3tech-oss/jwt-go"
)

const day = 24 * time.Hour

// ErrMissingCertificate from Auth0.
var ErrMissingCertificate = errors.New("missing certificate")

type jwksResponse struct {
	Keys []jsonWebKeys `json:"keys"`
}

type jsonWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type pem struct {
	cfg    *Config
	client *http.Client
}

func (p *pem) Certificate(ctx context.Context, token *jwt.Token) (string, error) {
	cert := ""

	httpReq, err := http.NewRequestWithContext(ctx, "GET", p.cfg.JSONWebKeySet, nil)
	if err != nil {
		return cert, err
	}

	httpResp, err := p.client.Do(httpReq)
	if err != nil {
		return cert, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return cert, ErrInvalidResponse
	}

	var resp jwksResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return cert, err
	}

	for k := range resp.Keys {
		if token.Header["kid"] == resp.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + resp.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		return cert, ErrMissingCertificate
	}

	return cert, nil
}

type cachedPEM struct {
	cfg   *Config
	cache *ristretto.Cache

	Certificator
}

func (p *cachedPEM) Certificate(ctx context.Context, token *jwt.Token) (string, error) {
	cacheKey := p.cfg.CacheKey("certificate")

	v, ok := p.cache.Get(cacheKey)
	if ok {
		return v.(string), nil
	}

	cert, err := p.Certificator.Certificate(ctx, token)
	if err != nil {
		return cert, err
	}

	p.cache.SetWithTTL(cacheKey, cert, 0, day)

	return cert, nil
}
