package crypto

import (
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
)

// Config is the root crypto configuration for a go-service based service.
//
// It composes configuration for multiple cryptographic subsystems that are implemented in subpackages
// under `github.com/alexfalkowski/go-service/v2/crypto/*`.
//
// # Optional pointers and "enabled" semantics
//
// All sub-config fields are pointers and are intentionally optional. A nil sub-config is generally
// interpreted as "disabled" by the corresponding subsystem (see each subpackage's `IsEnabled`
// convention where applicable).
//
// # Key material sources
//
// Many crypto subpackages support "source string" values (for example `env:NAME`, `file:/path`, or a
// raw literal) for secrets and key material. Those source strings are resolved by `os.FS.ReadSource`
// during construction or use, depending on the subsystem.
type Config struct {
	// AES configures AES key material used by AES-based primitives (for example encrypt/decrypt helpers).
	AES *aes.Config `yaml:"aes,omitempty" json:"aes,omitempty" toml:"aes,omitempty"`

	// Ed25519 configures Ed25519 public/private key material used for signing and verification.
	Ed25519 *ed25519.Config `yaml:"ed25519,omitempty" json:"ed25519,omitempty" toml:"ed25519,omitempty"`

	// HMAC configures HMAC key material used for message authentication.
	HMAC *hmac.Config `yaml:"hmac,omitempty" json:"hmac,omitempty" toml:"hmac,omitempty"`

	// RSA configures RSA public/private key material used for signing/verification and encryption/decryption
	// (depending on which helpers are wired by the service).
	RSA *rsa.Config `yaml:"rsa,omitempty" json:"rsa,omitempty" toml:"rsa,omitempty"`

	// SSH configures SSH public/private key material used for SSH signing and verification.
	SSH *ssh.Config `yaml:"ssh,omitempty" json:"ssh,omitempty" toml:"ssh,omitempty"`
}

// IsEnabled reports whether crypto configuration is enabled.
//
// By convention, a nil *Config is treated as "crypto disabled" by wiring that depends on this configuration.
func (c *Config) IsEnabled() bool {
	return c != nil
}
