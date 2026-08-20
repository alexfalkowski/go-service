// Package stream provides streaming (multi-value) encode/decode interfaces used by go-service.
//
// Unlike [github.com/alexfalkowski/go-service/v2/encoding], which models exactly one value per
// Encode/Decode call against a writer/reader supplied per call, this package models a sequence of
// values written to, or read from, one writer or reader bound at construction time. [Encoder] and
// [Decoder] are separate interfaces, rather than one combined interface, because the underlying
// per-kind codecs construct distinct encoder and decoder types that are each bound to a single
// direction.
//
// Per-kind implementations live in sibling packages
// (github.com/alexfalkowski/go-service/v2/encoding/stream/json, .../msgpack, .../gob, .../yaml), each
// exporting NewEncoder(io.Writer) [Encoder] and NewDecoder(io.Reader) [Decoder]. They wrap the same
// underlying libraries as the sibling encoding/<kind> packages, carrying over the same policy with two
// exceptions needed for multi-value streams: no output indentation (which would break newline-delimited
// framing) and no trailing-data rejection (a stream is expected to contain many values). Unknown
// fields are rejected by default and can be discarded at construction with
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithDiscardUnknown].
//
// Start with [Encoder] and [Decoder].
package stream
