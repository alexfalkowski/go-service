// Package msgpack provides a streaming MessagePack
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Encoder]/
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Decoder] pair used by go-service.
//
// It wraps [github.com/Basekick-Labs/msgpack/v6]'s own NewEncoder/NewDecoder types, which are each
// bound to one writer/reader for many Encode/Decode calls, unlike
// [github.com/alexfalkowski/go-service/v2/encoding/msgpack]'s Decode, which rejects trailing values
// because a stream is expected to contain many values.
//
// Start with the package-level constructors.
package msgpack
