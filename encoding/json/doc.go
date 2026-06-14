// Package json provides the go-service JSON import path.
//
// The package serves two related purposes:
//
//   - it exposes [Encoder], a thin adapter that satisfies the repository's
//     generic encoding abstraction while using readable indented encoding and
//     strict standard library decoding
//   - it re-exports common JSON helpers and types from the standard library so
//     packages in this repository can depend on a single go-service JSON import
//     path instead of importing encoding/json directly
//
// Marshal and Valid preserve the standard library's encoding/json semantics.
// Encoder and Unmarshal use the standard decoder with unknown fields and
// trailing values rejected.
package json
