// Package errors provides crypto-specific error values shared across go-service
// crypto helpers.
//
// The package centralizes sentinel errors that callers can compare with
// `errors.Is` regardless of which concrete crypto helper produced them.
//
// Current sentinels include:
//   - ErrInvalidMatch: a verification operation failed because the provided
//     signature, MAC, or hash did not match the message.
//   - ErrMissingKey: required key material was absent or resolved to empty bytes.
//   - ErrInvalidKeyType: key material decoded successfully, but the contained
//     key algorithm did not match what the caller expected.
//   - ErrInvalidKeySize: key material decoded successfully, but the key size
//     does not satisfy the package policy for the operation.
//   - ErrInvalidKeyFormat: key material decoded successfully, but used an
//     unsupported syntactic form.
//
// These sentinels are intentionally small and reusable so higher-level
// packages can add context by wrapping them without losing comparability.
package errors
