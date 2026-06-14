// Package compress provides compression abstractions and DI wiring for go-service.
//
// This package defines a small interface ([Compressor]) for compressing and
// decompressing `[]byte` payloads, plus a registry ([Map]) used to select an
// implementation by kind at runtime.
//
// The size passed to [Compressor.Compress] limits the uncompressed input. The
// size passed to [Compressor.Decompress] limits the decompressed output. When a
// size limit is exceeded, implementations return
// [github.com/alexfalkowski/go-service/v2/compress/errors.ErrTooLarge].
//
// # Registry
//
// [Map] is a simple kind-to-Compressor lookup (e.g. "zstd", "s2", "snappy", "none"). Callers typically
// obtain a *[Map] via DI and then use [Map.Get] to select the configured algorithm, falling back to "none"
// when the configured kind is not present.
//
// # Wiring
//
// [Module] wires the default compressor implementations and provides a *[Map] pre-populated with:
//
//   - "zstd"
//   - "s2"
//   - "snappy"
//   - "none"
//
// You can extend or override supported kinds by calling [Map.Register] after construction.
//
// Start with [Compressor], [Map], and [Module].
package compress
