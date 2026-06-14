package keys

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/token/errors"
)

// Map maps key ids to Ed25519 token key material.
type Map map[string]*Config

// Get returns the key with id.
func (m Map) Get(id string) *Config {
	if m == nil {
		return nil
	}

	return m[id]
}

// Config configures named Ed25519 token key material.
type Config struct {
	// Config contains the Ed25519 public/private key sources.
	*ed25519.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// Signer loads an Ed25519 signer for k.
//
// It returns [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig]
// when k is nil, its embedded key config is nil, or decoder is nil.
func (k *Config) Signer(decoder *pem.Decoder) (*ed25519.Signer, error) {
	if k == nil || k.Config == nil || decoder == nil {
		return nil, errors.ErrInvalidConfig
	}

	return ed25519.NewSigner(decoder, k.Config)
}

// Verifier loads an Ed25519 verifier for k.
//
// It returns [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig]
// when k is nil, its embedded key config is nil, or decoder is nil.
func (k *Config) Verifier(decoder *pem.Decoder) (*ed25519.Verifier, error) {
	if k == nil || k.Config == nil || decoder == nil {
		return nil, errors.ErrInvalidConfig
	}

	return ed25519.NewVerifier(decoder, k.Config)
}
