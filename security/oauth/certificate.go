package oauth

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/golang-jwt/jwt/v4"
)

const day = 24 * time.Hour

// ErrMissingCertificate from OAuth.
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

type certificate struct {
	cfg    *Config
	client *http.Client
}

func (c *certificate) Certificate(ctx context.Context, token *jwt.Token) (crypto.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.cfg.JSONWebKeySet, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, ErrInvalidResponse
	}

	var resp jwksResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	for k := range resp.Keys {
		if token.Header["kid"] == resp.Keys[k].Kid {
			return c.publicKey(resp.Keys[k].X5c[0])
		}
	}

	return nil, ErrMissingCertificate
}

func (c *certificate) publicKey(pk string) (ed25519.PublicKey, error) {
	return base64.StdEncoding.DecodeString(pk)
}

type cachedCertificate struct {
	cfg   *Config
	cache *ristretto.Cache

	Certificator
}

func (c *cachedCertificate) Certificate(ctx context.Context, token *jwt.Token) (crypto.PublicKey, error) {
	cacheKey := c.cfg.CacheKey("certificate")

	v, ok := c.cache.Get(cacheKey)
	if ok {
		return v.(crypto.PublicKey), nil
	}

	cert, err := c.Certificator.Certificate(ctx, token)
	if err != nil {
		return cert, err
	}

	c.cache.SetWithTTL(cacheKey, cert, 0, day)

	return cert, nil
}
