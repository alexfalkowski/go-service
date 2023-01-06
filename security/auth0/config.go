package auth0

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Config for Auth0.
type Config struct {
	URL           string `yaml:"url" json:"url"`
	ClientID      string `yaml:"client_id" json:"client_id"`
	ClientSecret  string `yaml:"client_secret" json:"client_secret"`
	Audience      string `yaml:"audience" json:"audience"`
	Issuer        string `yaml:"issuer" json:"issuer"`
	Algorithm     string `yaml:"algorithm" json:"algorithm"`
	JSONWebKeySet string `yaml:"json_web_key_set" json:"json_web_key_set"`
}

// CacheKey for config.
func (c *Config) CacheKey(prefix string) string {
	key := fmt.Sprintf("%s:%s:%s:%s:%s:%s", prefix, c.URL, c.ClientID, c.ClientSecret, c.Audience, c.JSONWebKeySet)
	h := sha256.New()

	h.Write([]byte(key))
	h.Sum(nil)

	return hex.EncodeToString(h.Sum(nil))
}
