// Package errors provides crypto-specific error values shared across go-service
// crypto helpers.
//
// The package centralizes sentinel errors that callers can compare with
// `errors.Is` regardless of which concrete crypto helper produced them.
//
// Current sentinels include:
//   - ErrInvalidMatch: a verification operation failed because the provided
//     signature, MAC, or hash did not match the message.
//   - ErrInvalidKeyType: key material decoded successfully, but the contained
//     key algorithm did not match what the caller expected.
//
// These sentinels are intentionally small and reusable so higher-level
// packages can add context by wrapping them without losing comparability.
package errors
