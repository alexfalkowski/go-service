// Package jwt provides JSON Web Token (JWT) issuance and verification for go-service.
//
// This package implements a JWT token kind intended for service-to-service and
// user/session authentication flows where a compact, self-contained token is
// required.
//
// # Algorithms and keys
//
// Tokens are signed using Ed25519 with the JWT "EdDSA" signing method
// (jwt.SigningMethodEdDSA). The implementation expects:
//
//   - an Ed25519 signing key for issuance, and
//   - an Ed25519 verification key for validation.
//
// The signing and verification keys are provided to NewToken by the caller (typically
// via DI wiring) using go-service crypto/ed25519 helpers.
//
// # Claims and headers
//
// Tokens use standard registered claims (jwt.RegisteredClaims). Issued tokens set:
//
//   - iss (issuer): from Config.Issuer
//   - aud (audience): from the Generate aud argument
//   - sub (subject): from the Generate sub argument
//   - iat (issued at): set to the current time
//   - nbf (not before): set to the current time
//   - exp (expiration): current time + Config.Expiration
//   - jti (JWT ID): generated via the provided id.Generator
//
// In addition, issued tokens set the JWT header:
//
//   - kid (key id): from Config.KeyID
//
// # Key ID (kid) enforcement
//
// Verification is intentionally strict about the "kid" header:
//
//   - The header must be present and non-empty.
//   - The value must exactly match Config.KeyID.
//
// This repository uses "kid" as part of the verification contract to prevent
// accepting tokens minted for a different key identity. If you mint test tokens
// using a third-party JWT library directly, ensure you set "kid" or verification
// will fail.
//
// # Verification semantics and errors
//
// Verify validates a token for a given audience and returns the subject (sub) on success.
// Verification enforces, in order:
//
//   - The signature algorithm is EdDSA.
//   - The "kid" header exists and matches the configured KeyID.
//   - The issuer claim matches the configured Issuer.
//   - The audience claim contains the expected audience.
//   - Registered claim validity (exp/nbf/iat) using jwt.RegisteredClaims.Valid.
//
// On failures, this package may return sentinel errors from token/errors for common
// classes of validation issues (for example ErrInvalidAlgorithm, ErrInvalidKeyID,
// ErrInvalidIssuer, ErrInvalidAudience). Other failures may be returned as-is from
// the upstream github.com/golang-jwt/jwt/v4 parser/validator.
//
// Note: Error wrapping and exact error values for parse/validation failures that are
// not explicitly mapped to token/errors are governed by the upstream JWT library.
//
// # Configuration and enablement
//
// Config provides the issuer, expiration, and key id settings used for issuance and
// verification. A nil *Config is treated as disabled: NewToken returns nil.
//
// Expiration is a Go duration string (time.ParseDuration format, such as "15m" or
// "24h"). Issuance uses MustParseDuration and will panic if Expiration is invalid.
// This is intended for strict startup/configuration paths; validate configuration
// earlier if you need non-panicking behavior.
//
// # Relationship to the top-level token facade
//
// Services often use the top-level token.Token facade (package token) which delegates
// to this implementation when Config.Kind == "jwt". This package focuses only on the
// JWT-specific format and validation rules.
package jwt
