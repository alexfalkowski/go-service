package security

import (
	"github.com/alexfalkowski/go-service/security/auth0"
)

// Config for security.
type Config struct {
	Auth0 auth0.Config `yaml:"auth0"`
}
