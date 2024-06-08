package crypto

import (
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	"github.com/alexfalkowski/go-service/crypto/ssh"
)

// IsEnabled for crypto.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for crypto.
type Config struct {
	AES     *aes.Config     `yaml:"aes,omitempty" json:"aes,omitempty" toml:"aes,omitempty"`
	Ed25519 *ed25519.Config `yaml:"ed25519,omitempty" json:"ed25519,omitempty" toml:"ed25519,omitempty"`
	HMAC    *hmac.Config    `yaml:"hmac,omitempty" json:"hmac,omitempty" toml:"hmac,omitempty"`
	RSA     *rsa.Config     `yaml:"rsa,omitempty" json:"rsa,omitempty" toml:"rsa,omitempty"`
	SSH     *ssh.Config     `yaml:"ssh,omitempty" json:"ssh,omitempty" toml:"ssh,omitempty"`
}
