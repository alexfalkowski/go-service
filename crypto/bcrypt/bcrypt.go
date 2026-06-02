package bcrypt

import (
	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/errors"
	"golang.org/x/crypto/bcrypt"
)

// DefaultCost is the cost used by Sign.
const DefaultCost = bcrypt.DefaultCost

// NewSigner constructs a bcrypt-based Signer intended for password hashing.
//
// The returned Signer uses [golang.org/x/crypto/bcrypt.DefaultCost] when hashing. If you need a different cost or more
// control over parameters, use [golang.org/x/crypto/bcrypt] directly.
func NewSigner() *Signer {
	return &Signer{}
}

// Cost returns the hashing cost used to create the provided bcrypt hash.
func Cost(hash []byte) (int, error) {
	return bcrypt.Cost(hash)
}

// Signer hashes and verifies secrets using bcrypt.
//
// This type is intended for password hashing and password hash verification.
// It is not a general-purpose message signing primitive.
//
// Sign returns a bcrypt hash for the provided secret (typically a password).
// Verify checks whether the provided bcrypt hash matches the provided secret.
type Signer struct{}

// Sign hashes msg using bcrypt with [golang.org/x/crypto/bcrypt.DefaultCost].
//
// The returned value is a bcrypt hash suitable for storage. The input msg is typically a user password.
// Callers should store only the returned hash, never the plaintext secret.
//
// This is a thin wrapper around [golang.org/x/crypto/bcrypt.GenerateFromPassword].
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(msg, DefaultCost)
}

// Verify checks that sig is a valid bcrypt hash for msg.
//
// This is a thin wrapper around [golang.org/x/crypto/bcrypt.CompareHashAndPassword].
// It returns nil if the hash matches, otherwise it returns [github.com/alexfalkowski/go-service/v2/crypto/errors.ErrInvalidMatch]
// joined with the bcrypt package error.
func (s *Signer) Verify(sig, msg []byte) error {
	err := bcrypt.CompareHashAndPassword(sig, msg)
	if err != nil {
		return errors.Join(crypto.ErrInvalidMatch, err)
	}

	return nil
}
