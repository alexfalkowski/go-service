package ed25519

import "crypto/ed25519"

// NewSigner for ed25519.
func NewSigner(cfg *Config) (*Signer, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	pri, err := cfg.PrivateKey()
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
