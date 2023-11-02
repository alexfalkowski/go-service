package security

import (
	"github.com/alexfalkowski/go-service/security/oauth"
)

// Config for security.
type Config struct {
	OAuth oauth.Config `yaml:"oauth" json:"oauth" toml:"oauth"`
}
