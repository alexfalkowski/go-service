package auth0

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Config for Auth0.
type Config struct {
	URL           string `envconfig:"AUTH0_URL" required:"true"`
	ClientID      string `envconfig:"AUTH0_CLIENT_ID" required:"true"`
	ClientSecret  string `envconfig:"AUTH0_CLIENT_SECRET" required:"true"`
	Audience      string `envconfig:"AUTH0_AUDIENCE" required:"true"`
	Issuer        string `envconfig:"AUTH0_ISSUER" required:"true"`
	JSONWebKeySet string `envconfig:"AUTH0_JSON_WEB_KEY_SET" required:"true"`
}

// CacheKey for config.
func (c *Config) CacheKey(suffix string) string {
	key := fmt.Sprintf("%s:%s:%s:%s:%s:%s", c.URL, c.ClientID, c.ClientSecret, c.Audience, c.JSONWebKeySet, suffix)
	h := sha256.New()

	h.Write([]byte(key))
	h.Sum(nil)

	return hex.EncodeToString(h.Sum(nil))
}
