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
//   - an active key id for issuance, and
//   - a named Ed25519 key set for issuance and validation.
//
// Key material is loaded from [Config.Keys] using the *[os.FS] passed to NewToken.
// Generate signs with [Config.Key], writes that id to the PASETO footer as
// "kid", and Verify selects the matching configured public key before signature
// verification.
//
// # Claims and footer
//
// Tokens are issued with common PASETO claims and standard identity fields.
// Issued tokens set:
//
//   - jti (token ID): generated via the provided id.Generator
//   - iat (issued at): set to the current time
//   - nbf (not before): set to the current time
//   - exp (expiration): current time + [Config.Expiration]
//   - iss (issuer): from [Config.Issuer]
//   - aud (audience): from the Generate aud argument
//   - sub (subject): from the Generate sub argument
//
// Issued tokens also carry a footer:
//
//   - kid (key id): from [Config.Key]
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
//   - for the expected audience (aud).
//
// The signature is verified using the Ed25519 public key selected by the footer
// "kid" value.
// After parsing and signature verification, go-service validates exp/nbf/iat
// using [Config.Leeway] as optional clock-skew tolerance and checks that the
// signed lifetime from iat to exp does not exceed [Config.Expiration].
//
// On failure, parser, rule, signature, and key-construction errors may come
// from the upstream PASETO library. Local config, subject, and signed-lifetime
// checks return shared sentinel errors from token/errors.
//
// # Configuration and enablement
//
// Config provides issuer, active key id, named key set, expiration, and optional
// leeway settings. A nil *[Config] is treated as disabled: NewToken returns nil.
//
// Expiration is a typed duration. In config files it is encoded using the standard
// Go duration string format (such as "15m" or "24h"), so invalid values fail during
// decoding rather than during token issuance.
//
// # Relationship to the top-level token facade
//
// Services often use the top-level [github.com/alexfalkowski/go-service/v2/token.Token] facade (package token) which
// delegates to this implementation when [Config.Kind] == "paseto". This package
// focuses only on the PASETO-specific format and validation rules.
package paseto
