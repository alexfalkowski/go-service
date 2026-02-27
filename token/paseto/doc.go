// Package paseto provides PASETO token generation and verification for go-service.
//
// This package implements a PASETO token kind intended for authentication and
// authorization flows where a compact, self-contained token is required.
//
// PASETO (Platform-Agnostic Security Tokens) is an alternative to JWT designed
// to avoid algorithm-downgrade pitfalls by versioning and clearer cryptographic
// choices. This package uses PASETO v4 public (asymmetric) tokens via the
// upstream implementation:
//
//	aidanwoods.dev/go-paseto
//
// # Algorithms and keys
//
// Tokens are signed and verified using PASETO v4 public tokens, which use Ed25519
// public-key signatures.
//
// The implementation expects:
//   - an Ed25519 signing key for issuance, and
//   - an Ed25519 verification key for validation.
//
// These keys are provided to NewToken by the caller (typically via DI wiring)
// using go-service crypto/ed25519 helpers.
//
// # Claims
//
// Tokens are issued with common PASETO claims and standard identity fields.
// Issued tokens set:
//
//   - jti (token ID): generated via the provided id.Generator
//   - iat (issued at): set to the current time
//   - nbf (not before): set to the current time
//   - exp (expiration): current time + Config.Expiration
//   - iss (issuer): from Config.Issuer
//   - aud (audience): from the Generate aud argument
//   - sub (subject): from the Generate sub argument
//
// # Verification semantics and errors
//
// Verify validates a token for a given audience and returns the subject (sub) on
// success.
//
// Verification is implemented by constructing a PASETO parser and applying a
// set of rules that typically include:
//
//   - issued by the configured issuer (iss),
//   - not expired at the time of verification,
//   - valid at the current time (iat/nbf window semantics as defined by the
//     upstream library),
//   - for the expected audience (aud).
//
// The signature is verified using the configured Ed25519 public key.
//
// On failure, this package returns errors from the upstream PASETO library or
// from key construction. It does not currently map failures onto the shared
// sentinel errors in token/errors.
//
// # Configuration and enablement
//
// Config provides issuer and expiration settings. A nil *Config is treated as
// disabled: NewToken returns nil.
//
// Expiration is a Go duration string (time.ParseDuration format, such as "15m" or
// "24h"). Issuance uses MustParseDuration and will panic if Expiration is invalid.
// This is intended for strict startup/configuration paths; validate configuration
// earlier if you need non-panicking behavior.
//
// # Secret field note
//
// Config contains a Secret field described as PASETO key material using the
// go-service “source string” convention (env:/file:/literal). However, the
// current implementation constructs PASETO v4 public tokens using Ed25519 key
// material supplied via crypto/ed25519.Signer and crypto/ed25519.Verifier passed
// to NewToken. The Secret field is therefore not consumed directly by this
// package’s Token implementation.
//
// If your service wants to source key material from config, resolve Config.Secret
// (using os.FS.ReadSource) and construct the Ed25519 signer/verifier accordingly
// in your wiring layer.
//
// # Relationship to the top-level token facade
//
// Services often use the top-level token.Token facade (package token) which
// delegates to this implementation when Config.Kind == "paseto". This package
// focuses only on the PASETO-specific format and validation rules.
package paseto
