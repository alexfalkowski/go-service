// Package ulid provides ULID generation and parsing helpers used by go-service.
//
// This package integrates ULID generation behind the go-service ID abstraction.
// ULID values are operational identifiers with time-sortable characteristics;
// they are not secret material or bearer tokens.
//
// Use [NewGenerator] to generate ULID strings and [Parse] to decode them into a [ULID].
package ulid
