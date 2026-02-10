// Package meta provides context-scoped metadata storage and helpers for go-service.
//
// This package provides a small key-value attribute store backed by context.Context. Values are stored as meta.Value
// so callers can control how attributes are rendered (normal, blank, ignored, or redacted) when exporting metadata
// to strings for logging and transport propagation.
//
// Start with `WithAttribute` / `Attribute` for working with arbitrary attributes and `Strings` / `SnakeStrings` /
// `CamelStrings` for exporting attributes as string maps.
package meta
