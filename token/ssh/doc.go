// Package ssh provides an SSH-style token format for go-service.
//
// This package implements a simple signed token scheme using SSH public key
// cryptography. It is intentionally different from claim-based tokens (JWT/PASETO):
// it does not encode audiences, issuers, expiration, or other claims. Instead, it
// provides a compact token that binds a key name to a signature.
//
// # Token format
//
// Tokens are ASCII strings of the form:
//
//	<name>-<base64(signature)>
//
// Where:
//
//   - <name> is the logical name of the signing key (for example "primary").
//   - signature is produced by signing the bytes of <name> with the configured SSH
//     private key.
//   - base64(signature) is the standard base64 encoding of the raw signature bytes.
//
// The separator is the first "-" in the token string. Anything after the first "-"
// is treated as the base64-encoded signature.
//
// # Signing keys and verification keys
//
// Configuration is provided via Config:
//
//   - Config.Key is the single signing key used for Generate.
//   - Config.Keys is a set of named public keys used for Verify.
//
// Verification is “name-based”: Verify extracts <name> from the token and then
// looks up a matching public key configuration in Config.Keys (via Keys.Get(name)).
// If no key with that name exists, verification fails.
//
// This design supports key rotation and multi-key verification: you can mint tokens
// with the active signing key name while allowing verification against multiple
// historical/active public keys by including them in Config.Keys.
//
// # Key material loading and “source strings”
//
// The Token constructor accepts an *os.FS and uses go-service crypto/ssh helpers to
// load key material based on the embedded crypto/ssh.Config in Key.
//
// Those configs commonly support go-service “source strings” for key sources
// (for example env:/file:/literal). Resolution and filesystem behavior depend on
// the go-service os.FS and crypto/ssh packages used by your wiring.
//
// # Error handling expectations
//
// Verify returns the key name on success (the <name> prefix from the token).
// On failure, it returns an error. Common failure modes include:
//
//   - token does not contain the "-" separator,
//   - no verification key exists for the extracted name,
//   - base64 decoding fails,
//   - signature verification fails,
//   - key material cannot be loaded.
//
// Some invalid-token cases are intentionally collapsed into a generic “invalid match”
// class so callers do not learn whether a name exists or which specific check failed.
// Callers that need fine-grained diagnostics should add logging/metrics at the call
// site rather than relying on error text.
//
// # Security notes
//
// This scheme authenticates possession of a key (via signature verification) and binds
// that to a logical name. It does not provide expiration or replay protection by
// itself. If your use case requires time-bounded validity, nonce/jti semantics, or
// audience restrictions, prefer JWT or PASETO token kinds or layer additional checks
// at a higher level.
//
// # Relationship to the top-level token facade
//
// Services often use the top-level token.Token facade (package token), which delegates
// to this implementation when Config.Kind == "ssh".
package ssh
