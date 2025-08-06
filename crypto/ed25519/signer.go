package ed25519

import (
	"crypto/ed25519"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
)

// NewSigner for ed25519.
func NewSigner(decoder *pem.Decoder, cfg *Config) (*Signer, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	pri, err := cfg.PrivateKey(decoder)
	if err != nil {
		return nil, err
	}

	return &Signer{PrivateKey: pri}, nil
}

// Signer for ed25519.
type Signer struct {
	PrivateKey ed25519.PrivateKey
}

// Sign for ed25519.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(s.PrivateKey, msg), nil
}
