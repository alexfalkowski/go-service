package ssh

import (
	"os"
)

// IsEnabled for ssh.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (
	// PublicKey for ssh.
	PublicKey string

	// PrivateKey for ssh.
	PrivateKey string

	// Config for ssh.
	Config struct {
		Public  PublicKey  `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
		Private PrivateKey `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
	}
)

// GetPublic for ssh.
func (c *Config) GetPublic() ([]byte, error) {
	return os.ReadFile(string(c.Public))
}

// GetPrivate for ssh.
func (c *Config) GetPrivate() ([]byte, error) {
	return os.ReadFile(string(c.Private))
}
