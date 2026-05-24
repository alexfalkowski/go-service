package ssh

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/slices"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Config configures the SSH-style token implementation.
//
// This token kind uses a simple signed token format (see package ssh docs) and requires
// SSH key material for:
//   - signing (minting tokens), and
//   - verification (validating tokens).
//
// Config separates those concerns:
//
//   - Key is the single signing key used by Token.Generate.
//   - Keys is the set of verification keys that Token.Verify may use.
//   - Expiration is how long newly generated tokens are valid.
//
// # Key rotation and multi-key verification
//
// Verification is name-based: the signed token claims embed a key id, and
// verification selects a matching public key config from Keys (via Keys.Get(name)).
// This design supports key rotation by allowing you to:
//
//   - mint new tokens with the active signing key name, and
//   - continue verifying older tokens by keeping historical public keys in Keys.
//
// Note: This package does not enforce that Key.Name exists in Keys. If you want tokens
// minted by Key to be verifiable by this same Config, include the corresponding public
// key entry in Keys under the same name.
//
// # Enablement
//
// Enablement is modeled by presence and content: a nil *Config is disabled, and a config
// with neither Key nor Keys is disabled (see IsEnabled).
type Config struct {
	// Key is the signing key configuration used to mint SSH-style tokens.
	//
	// This should include the private key material (via the embedded crypto/ssh.Config
	// fields) and a logical Name that will be embedded in minted tokens.
	//
	// If Key is nil, Token.Generate will not be usable.
	Key *Key `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`

	// Keys is the set of verification keys that may be used to validate SSH-style tokens.
	//
	// Verification uses the token's embedded key id to select a key from this set. If Keys
	// is empty or does not contain the token's key id, verification fails.
	//
	// If Keys is nil/empty, Token.Verify will not be usable.
	Keys Keys `yaml:"keys,omitempty" json:"keys,omitempty" toml:"keys,omitempty"`

	// Expiration is the duration used to set token expiration.
	//
	// In config files it is encoded as a Go duration string, for example "15m" or "1h".
	Expiration time.Duration `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty" validate:"gt=0"`
}

// IsEnabled reports whether SSH token configuration is enabled.
//
// It returns true when the receiver is non-nil and at least one of Key or Keys is present.
func (c *Config) IsEnabled() bool {
	return c != nil && (c.Key != nil || c.Keys != nil)
}

// Key describes SSH key material configuration along with its logical name.
//
// The embedded crypto/ssh.Config provides the public/private key source configuration
// used by go-service crypto/ssh helpers (typically via an os.FS).
//
// The Name identifies the key logically and is used to select keys during verification.
type Key struct {
	// Config contains the SSH key material configuration (public/private key sources).
	*ssh.Config `yaml:",inline" json:",inline" toml:",inline"`

	// Name is the logical key name used to select a key (for example via Keys.Get).
	//
	// For signing, this name is embedded into the signed token claims as the key id.
	// For verification, this name is used as the lookup key into Keys.
	Name string `yaml:"name,omitempty" json:"name,omitempty" toml:"name,omitempty"`
}

// Keys is a list of named SSH keys.
//
// This is used for verification key selection. Names are expected (but not required)
// to be unique; Get returns the first match.
type Keys []*Key

// Get returns the key with the given name, or nil if no matching key exists.
//
// If multiple keys share the same Name, Get returns the first match.
// Nil entries in the slice are ignored.
func (c Keys) Get(name string) *Key {
	key, _ := slices.ElemFunc(c, func(k *Key) bool { return k != nil && k.Name == name })

	return key
}
