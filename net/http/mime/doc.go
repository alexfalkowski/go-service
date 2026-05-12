// Package mime defines common MIME media type constants used by go-service.
//
// This package centralizes the media types used across go-service HTTP components for consistent
// Content-Type negotiation and response encoding.
//
// The constants in this package represent full media type strings (sometimes including charset
// parameters), suitable for use in HTTP headers such as:
//
//	Content-Type: application/json
//
// or:
//
//	Content-Type: text/plain; charset=utf-8
//
// Start with the `*MediaType` constants (e.g. JSONMediaType, ProtobufMediaType, TextMediaType).
package mime
