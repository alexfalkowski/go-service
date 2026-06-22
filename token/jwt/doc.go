// Package jwt provides JSON Web Token (JWT) issuance and verification for go-service.
//
// This package implements a JWT token kind intended for service-to-service and
// user/session authentication flows where a compact, self-contained token is
// required.
//
// # Algorithms and keys
//
// Tokens are signed using Ed25519 with the JWT "EdDSA" signing method
// ([github.com/golang-jwt/jwt/v4.SigningMethodEdDSA]). The implementation expects:
//
//   - an active key id for issuance, and
//   - a named Ed25519 key set for issuance and validation.
//
// Key material is loaded from [Config.Keys] using the *[os.FS] passed to NewToken.
// Generate signs with [Config.Key], while Verify reads the token's "kid" header
// and selects the matching configured public key.
//
// # Claims and headers
//
// Tokens use standard registered claims ([github.com/golang-jwt/jwt/v4.RegisteredClaims]). Issued tokens set:
//
//   - iss (issuer): from [Config.Issuer]
//   - aud (audience): from the Generate aud argument
//   - sub (subject): from the Generate sub argument
//   - iat (issued at): set to the current time
//   - nbf (not before): set to the current time
//   - exp (expiration): current time + [Config.Expiration]
//   - jti (JWT ID): generated via the provided id.Generator
//
// In addition, issued tokens set the JWT header:
//
//   - kid (key id): from [Config.Key]
//
// # Key ID (kid) enforcement
//
// Verification is intentionally strict about the "kid" header:
//
//   - The header must be present and non-empty.
//   - The value must select a configured entry from [Config.Keys].
//
// This repository uses "kid" as part of the verification contract to prevent
// accepting tokens minted for an untrusted key identity. If you mint test tokens
// using a third-party JWT library directly, ensure you set "kid" and configure
// the corresponding public key or verification will fail.
//
// # Verification semantics and errors
//
// Verify validates a token for a given audience and returns the subject (sub) on success.
// Verification enforces:
//
//   - The signature algorithm is EdDSA.
//   - The "kid" header exists and selects a configured verification key.
//   - The issuer claim matches the configured Issuer.
//   - The audience claim contains the expected audience.
//   - Registered claim time validity (exp/nbf/iat) using [Config.Leeway] as
//     optional clock-skew tolerance.
//   - The signed lifetime from iat to exp does not exceed [Config.Expiration].
//
// Upstream parsing and signature validation can fail before later go-service
// issuer, audience, time-validity, and lifetime checks run.
//
// On failures, this package may return sentinel errors from token/errors for common
// classes of validation issues (for example ErrInvalidAlgorithm, ErrInvalidKeyID,
// ErrInvalidIssuer, ErrInvalidAudience). Other failures may be returned as-is from
// the upstream [github.com/golang-jwt/jwt/v4] parser/validator.
//
// Note: Error wrapping and exact error values for parse/validation failures that are
// not explicitly mapped to token/errors are governed by the upstream JWT library.
//
// # Configuration and enablement
//
// Config provides the issuer, active key id, named key set, expiration, and optional
// leeway settings used for issuance and verification. A nil *[Config] is treated as disabled:
// NewToken returns nil.
//
// Expiration is a typed duration. In config files it is encoded using the standard
// Go duration string format (such as "15m" or "24h"), so invalid values fail during
// decoding rather than during token issuance.
//
// # Relationship to the top-level token facade
//
// Services often use the top-level [github.com/alexfalkowski/go-service/v2/token.Token] facade (package token) which delegates
// to this implementation when [Config.Kind] == "jwt". This package focuses only on the
// JWT-specific format and validation rules.
package jwt
