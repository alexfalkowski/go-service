package keys

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-sync"
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

	signer   sync.Pointer[ed25519.Signer]
	verifier sync.Pointer[ed25519.Verifier]
}

// Signer loads an Ed25519 signer for k.
//
// The signer is resolved at most once per process lifetime and reused for
// subsequent calls. Only a successful load is cached, so a transient read or
// parse failure is retried on the next call rather than being cached forever.
//
// It returns [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig]
// when k is nil, its embedded key config is nil, or decoder is nil.
func (k *Config) Signer(decoder *pem.Decoder) (*ed25519.Signer, error) {
	if k == nil || k.Config == nil || decoder == nil {
		return nil, errors.ErrInvalidConfig
	}

	if s := k.signer.Load(); s != nil {
		return s, nil
	}

	s, err := ed25519.NewSigner(decoder, k.Config)
	if err != nil {
		return nil, err
	}

	if !k.signer.CompareAndSwap(nil, s) {
		s = k.signer.Load()
	}

	return s, nil
}

// Verifier loads an Ed25519 verifier for k.
//
// The verifier is resolved at most once per process lifetime and reused for
// subsequent calls. Only a successful load is cached, so a transient read or
// parse failure is retried on the next call rather than being cached forever.
//
// It returns [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig]
// when k is nil, its embedded key config is nil, or decoder is nil.
func (k *Config) Verifier(decoder *pem.Decoder) (*ed25519.Verifier, error) {
	if k == nil || k.Config == nil || decoder == nil {
		return nil, errors.ErrInvalidConfig
	}

	if v := k.verifier.Load(); v != nil {
		return v, nil
	}

	v, err := ed25519.NewVerifier(decoder, k.Config)
	if err != nil {
		return nil, err
	}

	if !k.verifier.CompareAndSwap(nil, v) {
		v = k.verifier.Load()
	}

	return v, nil
}
