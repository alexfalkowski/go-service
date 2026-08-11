// Package compress provides zstd and gzip HTTP handler and transport wrappers.
//
// [GzipHandler] compresses server responses for clients that advertise zstd or gzip support, and [Transport]
// configures an outbound RoundTripper to advertise zstd and gzip support and transparently decompress responses.
// Set [HeaderNoCompression] on a response before it is written when an endpoint must remain uncompressed,
// such as a streaming response.
//
// The exported helpers are thin wrappers around the underlying implementation so callers can keep response
// compression integration on the go-service import path.
package compress
