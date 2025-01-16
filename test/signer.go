package test

import (
	"crypto/ed25519"
)

// BadEd25519Signer for test.
type BadEd25519Signer struct{}

func (s *BadEd25519Signer) PublicKey() ed25519.PublicKey {
	return nil
}

func (s *BadEd25519Signer) PrivateKey() ed25519.PrivateKey {
	return nil
}

func (s *BadEd25519Signer) Sign(_ string) (string, error) {
	return "", ErrInvalid
}

func (s *BadEd25519Signer) Verify(_, _ string) error {
	return ErrInvalid
}
