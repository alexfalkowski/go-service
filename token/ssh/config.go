package ssh

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	token "github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-sync"
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
//   - Key is the active signing key id used by [Token.Generate].
//   - Keys is the named key set that [Token.Generate] and [Token.Verify] may use.
//   - Expiration is how long newly generated tokens are valid.
//
// # Key rotation and multi-key verification
//
// Verification is id-based: the signed token claims embed a key id, and
// verification selects a matching public key config from Keys (via [Keys.Get]).
// This design supports key rotation by allowing you to:
//
//   - mint new tokens with the active signing key id, and
//   - continue verifying older tokens by keeping historical public keys in Keys.
//
// The active Key entry must exist in Keys and include private key material when
// generating tokens.
//
// # Enablement
//
// Enablement is modeled by presence and content: a nil *[Config] is disabled, and a config
// with neither Key nor Keys is disabled (see IsEnabled).
type Config struct {
	// Keys is the named key set used to mint and verify SSH-style tokens.
	//
	// Generation uses Key to select private key material. Verification uses the
	// token's embedded key id to select public key material.
	//
	// If Keys is nil/empty, Token.Verify will not be usable.
	Keys Keys `yaml:"keys,omitempty" json:"keys,omitempty" toml:"keys,omitempty"`

	// Key is the active signing key id used to mint SSH-style tokens.
	//
	// The selected key id is embedded in minted tokens as both the "kid" and "sub"
	// claims because SSH tokens authenticate the trusted peer key itself.
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`

	// Expiration is the duration used to set token expiration.
	//
	// In config files it is encoded as a Go duration string, for example "15m" or "1h".
	// It must use whole-second precision.
	Expiration time.Duration `yaml:"exp,omitempty" json:"exp,omitempty" toml:"exp,omitempty" validate:"duration_second_precision"`

	// Leeway is the optional clock-skew tolerance applied during verification.
	//
	// A zero value keeps strict time validation. Non-zero values allow iat to be
	// slightly in the future and exp to be slightly in the past while preserving
	// the signed lifetime cap enforced from iat to exp.
	Leeway time.Duration `yaml:"leeway,omitempty" json:"leeway,omitempty" toml:"leeway,omitempty" validate:"omitempty,duration_second_precision"`
}

// IsEnabled reports whether SSH token configuration is enabled.
//
// It returns true when the receiver is non-nil and at least one of Key or Keys is present.
func (c *Config) IsEnabled() bool {
	return c != nil && (!strings.IsEmpty(c.Key) || len(c.Keys) > 0)
}

// Key describes SSH key material configuration.
//
// The embedded [github.com/alexfalkowski/go-service/v2/crypto/ssh.Config] provides the public/private key source configuration
// used by go-service crypto/ssh helpers (typically via an [os.FS]).
type Key struct {
	// Config contains the SSH key material configuration (public/private key sources).
	*ssh.Config `yaml:",inline" json:",inline" toml:",inline"`

	signer   sync.Pointer[ssh.Signer]
	verifier sync.Pointer[ssh.Verifier]
}

// Signer loads an SSH signer for k.
//
// The signer is resolved at most once per process lifetime and reused for
// subsequent calls. Only a successful load is cached, so a transient read or
// parse failure is retried on the next call rather than being cached forever.
//
// It returns [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig]
// when k is nil, its embedded key config is nil, or fs is nil.
func (k *Key) Signer(fs *os.FS) (*ssh.Signer, error) {
	if k == nil || k.Config == nil || fs == nil {
		return nil, token.ErrInvalidConfig
	}

	if s := k.signer.Load(); s != nil {
		return s, nil
	}

	s, err := ssh.NewSigner(fs, k.Config)
	if err != nil {
		return nil, err
	}

	if !k.signer.CompareAndSwap(nil, s) {
		s = k.signer.Load()
	}

	return s, nil
}

// Verifier loads an SSH verifier for k.
//
// The verifier is resolved at most once per process lifetime and reused for
// subsequent calls. Only a successful load is cached, so a transient read or
// parse failure is retried on the next call rather than being cached forever.
//
// It returns [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig]
// when k is nil, its embedded key config is nil, or fs is nil.
func (k *Key) Verifier(fs *os.FS) (*ssh.Verifier, error) {
	if k == nil || k.Config == nil || fs == nil {
		return nil, token.ErrInvalidConfig
	}

	if v := k.verifier.Load(); v != nil {
		return v, nil
	}

	v, err := ssh.NewVerifier(fs, k.Config)
	if err != nil {
		return nil, err
	}

	if !k.verifier.CompareAndSwap(nil, v) {
		v = k.verifier.Load()
	}

	return v, nil
}

// Keys maps key ids to SSH key material.
//
// This is used for signing and verification key selection.
type Keys map[string]*Key

// Get returns the key with the given id, or nil if no matching key exists.
func (c Keys) Get(id string) *Key {
	if c == nil {
		return nil
	}

	return c[id]
}
