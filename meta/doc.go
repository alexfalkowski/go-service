// Package meta provides context-scoped metadata storage and helpers for go-service.
//
// This package implements a small attribute store backed by context.Context. Attributes are stored as
// key→meta.Value pairs, where Value carries both the underlying string and rendering semantics.
//
// # Storage model
//
// Attributes are stored on the context under a single internal key as a map-like Storage. Helpers such as
// WithAttribute return a derived context containing the updated storage.
//
// # Value rendering semantics
//
// Values are stored as meta.Value so callers can control how attributes are rendered when exporting them:
//   - normal: renders the underlying value as-is
//   - blank: represents "no value" and renders as empty
//   - ignored: stores the underlying value but renders as empty (useful to keep the value in-context while
//     preventing export to logs/headers)
//   - redacted: renders as a fixed-length mask (asterisks) while retaining the underlying value in-context
//
// # Export helpers
//
// Stored attributes can be exported to plain string maps for logging and transport propagation using:
//   - Strings: keys unchanged
//   - SnakeStrings: keys converted to snake_case
//   - CamelStrings: keys converted to lowerCamelCase
//
// Export helpers skip attributes whose rendered string is empty. A prefix may be prepended to each exported key.
//
// Start with `WithAttribute` / `Attribute` for arbitrary attributes, `Value` constructors (String/Blank/Ignored/Redacted)
// for controlling rendering, and `Strings` / `SnakeStrings` / `CamelStrings` for exporting attributes.
package meta
