// Package rsa provides RSA key loading, encryption, and decryption helpers for go-service.
//
// This package supports loading RSA keys from PEM-encoded PKCS#1 formats:
//
//   - public key: PEM block "RSA PUBLIC KEY" containing PKCS#1-encoded bytes (x509.ParsePKCS1PublicKey)
//   - private key: PEM block "RSA PRIVATE KEY" containing PKCS#1-encoded bytes (x509.ParsePKCS1PrivateKey)
//
// It also provides Encryptor/Decryptor helpers that use RSA-OAEP with SHA-512.
//
// Configuration values are typically loaded using the go-service "source string" pattern (for example "env:NAME",
// "file:/path", or a literal PEM value) via crypto/pem.Decoder.
package rsa
