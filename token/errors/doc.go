// Package errors defines shared sentinel errors used by go-service token implementations.
//
// This package centralizes common token validation failures so callers can:
//
//   - Check classes of failures consistently across token kinds (JWT, PASETO, SSH, etc.),
//     typically via errors.Is.
//   - Avoid importing implementation-specific packages just to compare error values.
//   - Map token failures to transport-appropriate responses (for example HTTP 401/403,
//     gRPC status codes, audit events, or metrics dimensions).
//
// # Sentinel errors and matching
//
// The exported variables in this package are intended to be used as sentinel errors.
// Token implementations should wrap these errors (for example using fmt.Errorf with %w
// or a helper that preserves an error chain) so callers can detect the failure reason
// while still preserving additional context.
//
// Example:
//
//	if errors.Is(err, tokenerrors.ErrInvalidAudience) {
//		// audience claim did not match expected audience
//	}
//
// # Meaning of common errors
//
// While concrete semantics are implementation-defined, the sentinels in this package
// generally correspond to:
//
//   - ErrInvalidMatch: cryptographic or structural mismatch (signature/MAC mismatch,
//     invalid encoding, etc.).
//   - ErrInvalidIssuer: issuer ("iss") claim mismatch.
//   - ErrInvalidAudience: audience ("aud") claim mismatch.
//   - ErrInvalidAlgorithm: algorithm mismatch (token signed/encrypted with an unexpected
//     algorithm).
//   - ErrInvalidKeyID: missing or unexpected key identifier (for example JWT "kid").
//   - ErrInvalidTime: time validity failure (expired, not-before in the future, etc.).
//
// # Notes
//
// This package does not provide error constructors or formatting helpers. If you need to
// attach more context (which claim failed, expected vs. actual, key name, etc.), wrap the
// sentinel error in the implementation package and preserve it in the error chain.
package errors
