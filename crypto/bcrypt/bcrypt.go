package bcrypt

import "golang.org/x/crypto/bcrypt"

// NewSigner returns a bcrypt-based Signer.
func NewSigner() *Signer {
	return &Signer{}
}

// Signer signs and verifies messages using bcrypt.
//
// Sign produces a bcrypt hash of the given message.
// Verify compares the given signature (bcrypt hash) against the given message.
type Signer struct{}

// Sign hashes msg using bcrypt with the default cost.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(msg, bcrypt.DefaultCost)
}

// Verify checks that sig is a valid bcrypt hash for msg.
func (s *Signer) Verify(sig, msg []byte) error {
	return bcrypt.CompareHashAndPassword(sig, msg)
}
