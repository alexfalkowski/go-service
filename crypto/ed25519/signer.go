package ed25519

import (
	"crypto/ed25519"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
)

// NewSigner constructs a Signer when configuration is enabled.
//
// If cfg is disabled, it returns (nil, nil). When enabled, it loads the private key using cfg.PrivateKey.
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

// Signer holds an Ed25519 private key used for signing messages.
type Signer struct {
	PrivateKey ed25519.PrivateKey
}

// Sign signs msg using Ed25519.
//
// This method does not fail as long as the private key is valid; it returns a nil error for API compatibility.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(s.PrivateKey, msg), nil
}
