// Package bcrypt provides bcrypt password hashing helpers for go-service.
//
// This package wraps golang.org/x/crypto/bcrypt to provide a simple Signer that:
//   - hashes passwords using bcrypt.GenerateFromPassword with bcrypt.DefaultCost, and
//   - verifies password hashes using bcrypt.CompareHashAndPassword.
//
// Start with `Signer` and `NewSigner`.
package bcrypt
