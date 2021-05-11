package auth0

import (
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/security/token"
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for Auth0.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// NewGenerator for Auth0.
func NewGenerator(cfg *Config, client *http.Client) token.Generator {
	return &generator{cfg: cfg, client: client}
}
