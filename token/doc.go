// Package token provides token generation and verification helpers used by go-service.
//
// This package defines common token concepts and shared helpers used by concrete
// token implementations (for example JWT, PASETO, and SSH).
//
// It also provides a small facade type (Token) that delegates token generation and
// verification to the configured implementation so callers can depend on a single
// entry point when the token kind is selected by configuration.
//
// # Supported token kinds
//
// The top-level Token facade supports the following kinds, selected by Config.Kind:
//
//   - "jwt": JSON Web Tokens signed using Ed25519 (see the token/jwt package).
//   - "paseto": PASETO v4 public tokens (see the token/paseto package).
//   - "ssh": SSH-style signed tokens (see the token/ssh package).
//
// Concrete packages document their own token formats, claims, cryptographic
// algorithms, and validation semantics.
//
// # Facade behavior and unknown kinds
//
// The Token facade is intentionally conservative when Config.Kind is unknown:
//
//   - Generate returns (nil, token/errors.ErrInvalidConfig).
//   - Verify returns (strings.Empty, token/errors.ErrInvalidConfig).
//
// This makes “unknown token kind” fail closed instead of behaving like “feature
// disabled” in wiring scenarios. Callers should treat ErrInvalidConfig as a
// startup or deployment configuration issue.
//
// # Configuration and enablement
//
// This package does not enforce that nested config blocks (JWT/Paseto/SSH) are
// present when the corresponding kind is selected. The concrete token constructors
// typically treat a nil *Config as disabled and may return nil implementations.
// Ensure your configuration is consistent with the selected kind.
package token
