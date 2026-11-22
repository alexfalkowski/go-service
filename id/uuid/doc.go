// Package uuid provides UUID generation and parsing helpers used by go-service.
//
// This package integrates UUID generation behind the go-service ID abstraction, using the
// standard library's uuid package. UUID values produced here are operational identifiers,
// not secret material or bearer tokens.
//
// Use [NewGenerator] to generate UUIDv7 strings and [Parse] to decode UUID text into a [UUID].
package uuid
