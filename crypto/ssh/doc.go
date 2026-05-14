// Package ssh provides Ed25519 SSH key loading and signing/verification helpers for go-service.
//
// This package supports:
//   - generating Ed25519 SSH key pairs,
//   - loading Ed25519 keys from SSH key formats, and
//   - constructing Signer/Verifier helpers around that key material.
//
// Supported key encodings are:
//
//   - public key: SSH authorized_keys format (parsed via ssh.ParseAuthorizedKey)
//   - private key: SSH private key format (parsed via ssh.ParseRawPrivateKey)
//
// Configuration values are typically loaded using the go-service "source string" pattern (for example "env:NAME",
// "file:/path", or a literal key value) via os.FS.ReadSource.
//
// The package treats a nil *Config as disabled in constructor-style helpers such
// as NewSigner and NewVerifier, allowing SSH signing/verification to be wired as
// an optional capability.
//
// If the provided key material is valid SSH key data but not an Ed25519 key,
// key parsing helpers return crypto/errors.ErrInvalidKeyType instead of
// panicking. Callers can use errors.Is(err, crypto.ErrInvalidKeyType) to
// detect that case while still receiving wrapped context about the actual
// decoded type.
package ssh
