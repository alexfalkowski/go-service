// Package ssh provides an SSH-style token format for go-service.
//
// This package implements a simple signed token scheme using SSH public key
// cryptography. It is intentionally different from issuer-based tokens
// (JWT/PASETO): it does not encode issuers or arbitrary claims. Instead, it
// provides a compact token that binds a key id, subject, audience, issued-at
// time, and expiration claims to a signature.
//
// # Token format
//
// Tokens are ASCII strings of the form:
//
//	<base64(json-claims)>.<base64(signature)>
//
// Where:
//
//   - json-claims contains "kid" (the logical signing key id), "sub" (the
//     subject, which must equal "kid"), and "aud" (the expected audience, such
//     as an HTTP path or gRPC method), plus "iat" and "exp" Unix nanosecond
//     timestamps.
//   - signature is produced by signing the exact JSON claims bytes with the
//     configured SSH private key.
//   - base64(signature) is the standard base64 encoding of the raw signature bytes.
//
// # Signing keys and verification keys
//
// Configuration is provided via Config:
//
//   - [Config.Key] is the active signing key id used for Generate.
//   - [Config.Keys] is a named key set used for Generate and Verify.
//   - [Config.Expiration] controls how long generated tokens remain valid.
//
// Verification is id-based: Verify extracts kid from the signed claims and
// then looks up a matching public key configuration in [Config.Keys] (via
// [Keys.Get](kid)). If no key with that id exists, verification fails. Verify
// also requires "sub" to match "kid" and returns that subject on success.
//
// This design supports key rotation and multi-key verification: you can mint tokens
// with the active signing key id while allowing verification against multiple
// historical/active public keys by including them in [Config.Keys].
//
// # Key material loading and "source strings"
//
// The Token constructor accepts an *[os.FS] and uses go-service crypto/ssh helpers to
// load key material based on the embedded crypto/ssh.Config selected from Keys.
//
// Those configs commonly support go-service "source strings" for key sources
// (for example "env:SSH_KEY", "file:/path/to/key", or a literal value).
// Resolution and filesystem behavior depend on the go-service [os.FS] and
// crypto/ssh packages used by your wiring.
//
// # Error handling expectations
//
// Verify returns the subject on success. For SSH tokens this subject must match
// the kid field from the signed claims. On failure, it returns an empty subject
// plus an error. Common failure modes include:
//
//   - token does not contain the "." separator,
//   - no verification key exists for the extracted id,
//   - the signed audience does not match the expected audience,
//   - the token is expired or not yet valid,
//   - base64 decoding fails,
//   - signature verification fails,
//   - key material cannot be loaded.
//
// Some invalid-token cases are intentionally collapsed into a generic "invalid match"
// class so callers do not learn whether a name exists or which specific check failed.
// Callers that need fine-grained diagnostics should add logging/metrics at the call
// site rather than relying on error text.
//
// # Security notes
//
// This scheme authenticates possession of a key (via signature verification) and binds
// that to a logical subject/key id, audience, and validity window. It does not provide
// nonce/jti replay protection for repeated calls to the same audience inside the
// validity window. If your use case requires one-time-use tokens, prefer JWT or
// PASETO token kinds with jti tracking or layer additional checks at a higher level.
//
// # Relationship to the top-level token facade
//
// Services often use the top-level [github.com/alexfalkowski/go-service/v2/token.Token] facade (package token), which delegates
// to this implementation when [Config.Kind] == "ssh".
package ssh
