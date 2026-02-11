package crypto

import (
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
)

// Config configures cryptographic primitives used by go-service.
//
// Individual sub-configs may be nil/disabled depending on which crypto features are enabled by the caller.
type Config struct {
	// AES configures AES key material used by AES-based primitives.
	AES *aes.Config `yaml:"aes,omitempty" json:"aes,omitempty" toml:"aes,omitempty"`

	// Ed25519 configures Ed25519 public/private key material used for signing and verification.
	Ed25519 *ed25519.Config `yaml:"ed25519,omitempty" json:"ed25519,omitempty" toml:"ed25519,omitempty"`

	// HMAC configures HMAC key material used for message authentication.
	HMAC *hmac.Config `yaml:"hmac,omitempty" json:"hmac,omitempty" toml:"hmac,omitempty"`

	// RSA configures RSA public/private key material used for signing/verification and related operations.
	RSA *rsa.Config `yaml:"rsa,omitempty" json:"rsa,omitempty" toml:"rsa,omitempty"`

	// SSH configures SSH public/private key material used for SSH-based crypto operations.
	SSH *ssh.Config `yaml:"ssh,omitempty" json:"ssh,omitempty" toml:"ssh,omitempty"`
}

// IsEnabled reports whether crypto configuration is enabled.
//
// A nil config is considered disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}
