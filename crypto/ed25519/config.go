package ed25519

import (
	"github.com/alexfalkowski/go-service/os"
)

// IsEnabled for ed25519.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (
	// PublicKey for ed25519.
	PublicKey string

	// PrivateKey  for ed25519.
	PrivateKey string

	// Config for ed25519.
	Config struct {
		Public  PublicKey  `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
		Private PrivateKey `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
	}
)

// GetPublic for ed25519.
func (c *Config) GetPublic() (string, error) {
	return os.ReadBase64File(string(c.Public))
}

// GetPrivate for ed25519.
func (c *Config) GetPrivate() (string, error) {
	return os.ReadBase64File(string(c.Private))
}
