// Package ed25519 provides Ed25519 key generation, signing, and verification helpers for go-service.
//
// This package supports:
//   - generating Ed25519 key pairs,
//   - loading Ed25519 keys from PEM-encoded X.509 containers, and
//   - constructing Signer/Verifier helpers around that key material.
//
// Supported key encodings are:
//   - public key: PEM block "PUBLIC KEY" containing PKIX-encoded bytes (x509.ParsePKIXPublicKey)
//   - private key: PEM block "PRIVATE KEY" containing PKCS#8-encoded bytes (x509.ParsePKCS8PrivateKey)
//
// Configuration values are typically loaded using the go-service "source string" pattern (for example "env:NAME",
// "file:/path", or a literal PEM value) via crypto/pem.Decoder.
//
// The package treats a nil *Config as disabled in constructor-style helpers such
// as NewSigner and NewVerifier, allowing callers to wire Ed25519 support as an
// optional dependency.
//
// If the decoded key material is valid but not Ed25519, key parsing helpers
// return crypto/errors.ErrInvalidKeyType instead of panicking. Callers can use
// errors.Is(err, crypto.ErrInvalidKeyType) to detect that case while still
// receiving wrapped context about the actual decoded type.
package ed25519
