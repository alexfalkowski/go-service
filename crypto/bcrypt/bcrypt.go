package bcrypt

import "golang.org/x/crypto/bcrypt"

// NewSigner constructs a bcrypt-based Signer intended for password hashing.
//
// The returned Signer uses bcrypt.DefaultCost when hashing. If you need a different cost or more
// control over parameters, use golang.org/x/crypto/bcrypt directly.
func NewSigner() *Signer {
	return &Signer{}
}

// Signer hashes and verifies secrets using bcrypt.
//
// This type is intended for password hashing and password hash verification.
// It is not a general-purpose message signing primitive.
//
// Sign returns a bcrypt hash for the provided secret (typically a password).
// Verify checks whether the provided bcrypt hash matches the provided secret.
type Signer struct{}

// Sign hashes msg using bcrypt with bcrypt.DefaultCost.
//
// The returned value is a bcrypt hash suitable for storage. The input msg is typically a user password.
// Callers should store only the returned hash, never the plaintext secret.
//
// This is a thin wrapper around bcrypt.GenerateFromPassword.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(msg, bcrypt.DefaultCost)
}

// Verify checks that sig is a valid bcrypt hash for msg.
//
// This is a thin wrapper around bcrypt.CompareHashAndPassword.
// It returns nil if the hash matches, otherwise it returns an error from the bcrypt package.
func (s *Signer) Verify(sig, msg []byte) error {
	return bcrypt.CompareHashAndPassword(sig, msg)
}
