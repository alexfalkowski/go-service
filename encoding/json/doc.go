// Package json provides the go-service JSON import path.
//
// The package serves two related purposes:
//
//   - it exposes [Encoder], a thin adapter that satisfies the repository's
//     generic encoding abstraction while preserving the default behavior of the
//     standard library encoder and decoder
//   - it re-exports common JSON helpers and types from the standard library so
//     packages in this repository can depend on a single go-service JSON import
//     path instead of importing encoding/json directly
//
// All marshaling and unmarshaling semantics remain those of the standard
// library's encoding/json package.
package json
