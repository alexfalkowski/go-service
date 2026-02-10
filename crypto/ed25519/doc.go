// Package ed25519 provides Ed25519 key generation, signing, and verification helpers for go-service.
//
// This package supports generating key pairs and loading keys from PEM-encoded X.509 formats:
//   - public key: PEM block "PUBLIC KEY" containing PKIX-encoded bytes (x509.ParsePKIXPublicKey)
//   - private key: PEM block "PRIVATE KEY" containing PKCS#8-encoded bytes (x509.ParsePKCS8PrivateKey)
//
// Configuration values are typically loaded using the go-service "source string" pattern (for example "env:NAME",
// "file:/path", or a literal PEM value) via crypto/pem.Decoder.
package ed25519
