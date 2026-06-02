// Package bcrypt provides bcrypt password hashing helpers for go-service.
//
// This package wraps [golang.org/x/crypto/bcrypt] to provide a simple Signer that:
//   - hashes passwords using [golang.org/x/crypto/bcrypt.GenerateFromPassword] with [golang.org/x/crypto/bcrypt.DefaultCost], and
//   - verifies password hashes using [golang.org/x/crypto/bcrypt.CompareHashAndPassword].
//
// Start with [Signer] and [NewSigner].
package bcrypt
