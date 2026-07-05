// Package content provides HTTP content negotiation helpers used by go-service.
//
// This package helps select an encoder/decoder based on HTTP media types (Content-Type and Accept) and
// provides small building blocks for content-aware request/response handling.
//
// # Media types and encoders
//
// The core type is [Content], which uses an `encoding.Map` registry to resolve an encoder by
// media subtype (e.g. "json", "hjson", "yaml", "toml", "proto").
//
// [Content] can derive a [Media] from either:
//   - an incoming HTTP request's Content-Type header, falling back to Accept ([Content.NewFromRequest]),
//   - an incoming HTTP request's Accept header, falling back to Content-Type ([Content.NewFromAccept]),
//   - an incoming HTTP request's Content-Type header ([Content.NewFromContentType]), or
//   - a raw media type string ([Content.NewFromMedia]).
//
// # Error payloads
//
// go-service uses a dedicated error media subtype ("error") to signal error payloads (typically
// rendered as plain text). When the subtype is "error", [NewMedia] returns a Media without an
// encoder and callers should treat the response body as an error message.
//
// # Defaults and fallbacks
//
// If media type parsing fails or the subtype is unknown, this package falls back to JSON.
//
// Start with [NewContent], [Content.NewFromRequest], and [Content.NewFromMedia].
package content
