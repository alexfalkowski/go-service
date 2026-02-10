// Package ssh provides Ed25519 SSH key loading and signing/verification helpers for go-service.
//
// This package supports generating and loading Ed25519 keys in SSH formats:
//
//   - public key: SSH authorized_keys format (parsed via ssh.ParseAuthorizedKey)
//   - private key: SSH private key format (parsed via ssh.ParseRawPrivateKey)
//
// Configuration values are typically loaded using the go-service "source string" pattern (for example "env:NAME",
// "file:/path", or a literal key value) via os.FS.ReadSource.
//
// Note: This package uses type assertions when parsing keys. If the provided key material is not an Ed25519 SSH key,
// key parsing helpers may panic due to type assertions.
package ssh
