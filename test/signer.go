package test

import (
	"crypto/ed25519"
)

// ErrEd25519Signer for test.
type ErrEd25519Signer struct{}

func (s *ErrEd25519Signer) PublicKey() ed25519.PublicKey {
	return nil
}

func (s *ErrEd25519Signer) PrivateKey() ed25519.PrivateKey {
	return nil
}

func (s *ErrEd25519Signer) Sign(_ string) (string, error) {
	return "", ErrInvalid
}

func (s *ErrEd25519Signer) Verify(_, _ string) error {
	return ErrInvalid
}
