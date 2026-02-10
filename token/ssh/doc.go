// Package ssh provides SSH-style token generation and verification for go-service.
//
// Tokens are generated as `<name>-<base64(signature)>`, where the signature is computed
// over the key name and verified using the configured public keys.
package ssh
