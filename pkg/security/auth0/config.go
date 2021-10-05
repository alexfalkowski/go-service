package auth0

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Config for Auth0.
type Config struct {
	URL           string `yaml:"url"`
	ClientID      string `yaml:"client_id"`
	ClientSecret  string `yaml:"client_secret"`
	Audience      string `yaml:"audience"`
	Issuer        string `yaml:"issuer"`
	Algorithm     string `yaml:"algorithm"`
	JSONWebKeySet string `yaml:"json_web_key_set"`
}

// CacheKey for config.
func (c *Config) CacheKey(prefix string) string {
	key := fmt.Sprintf("%s:%s:%s:%s:%s:%s", prefix, c.URL, c.ClientID, c.ClientSecret, c.Audience, c.JSONWebKeySet)
	h := sha256.New()

	h.Write([]byte(key))
	h.Sum(nil)

	return hex.EncodeToString(h.Sum(nil))
}
