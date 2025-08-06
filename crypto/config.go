package crypto

import (
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
)

// Config for crypto.
type Config struct {
	AES     *aes.Config     `yaml:"aes,omitempty" json:"aes,omitempty" toml:"aes,omitempty"`
	Ed25519 *ed25519.Config `yaml:"ed25519,omitempty" json:"ed25519,omitempty" toml:"ed25519,omitempty"`
	HMAC    *hmac.Config    `yaml:"hmac,omitempty" json:"hmac,omitempty" toml:"hmac,omitempty"`
	RSA     *rsa.Config     `yaml:"rsa,omitempty" json:"rsa,omitempty" toml:"rsa,omitempty"`
	SSH     *ssh.Config     `yaml:"ssh,omitempty" json:"ssh,omitempty" toml:"ssh,omitempty"`
}

// IsEnabled for crypto.
func (c *Config) IsEnabled() bool {
	return c != nil
}
