package ed25519

import (
	"crypto/ed25519"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
)

// NewSigner constructs an Ed25519 Signer when configuration is enabled.
//
// Disabled behavior: if cfg is nil (disabled), NewSigner returns (nil, nil).
//
// Enabled behavior: NewSigner loads and parses the Ed25519 private key using cfg.PrivateKey.
// It returns any error encountered during PEM decoding or key parsing.
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
	// PrivateKey is the Ed25519 private key used by Sign.
	PrivateKey ed25519.PrivateKey
}

// Sign signs msg using Ed25519 and returns the signature.
//
// Ed25519 signing does not return an error for a valid private key; this method returns a nil error
// for API compatibility with other signers.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(s.PrivateKey, msg), nil
}
