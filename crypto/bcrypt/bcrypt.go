package bcrypt

import "golang.org/x/crypto/bcrypt"

// NewSigner for bcrypt.
func NewSigner() *Signer {
	return &Signer{}
}

// Signer for bcrypt.
type Signer struct{}

// Sign for bcrypt.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(msg, bcrypt.DefaultCost)
}

// Verify for bcrypt.
func (s *Signer) Verify(sig, msg []byte) error {
	return bcrypt.CompareHashAndPassword(sig, msg)
}
