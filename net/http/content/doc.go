// Package content provides HTTP content negotiation helpers used by go-service.
//
// This package helps select an encoder/decoder based on HTTP media types (Content-Type and Accept) and
// provides small building blocks for content-aware request/response handling.
//
// # Media types and encoders
//
// The core type is [Content], which uses an `encoding.Map` registry to resolve an encoder by
// media subtype (e.g. "json", "hjson", "yaml", "toml", "protobuf"). A subtype absent from the registry,
// such as "proto" or "pbjson", may still resolve if this package aliases it to a registered kind (see
// [unaryKind]'s local alias table).
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
// For outbound negotiation ([Content.NewFromMedia] and Accept-based response encoder selection), if
// media type parsing fails or the subtype is unknown, this package falls back to JSON.
//
// Inbound request-body decoding ([Content.NewFromRequestBody] and its streaming counterpart
// [Content.NewStreamFromContentType]) narrows this: an absent Content-Type still defaults to JSON,
// since the caller asserted nothing, but an unparseable or unregistered Content-Type is rejected
// instead of silently decoded as a different format the caller did not send.
//
// # The decoder-bounds rule
//
// A wire format is admissible for decoding untrusted input — a server request body, or a streaming
// request value — only when its decoder is both ratio-bounded (it never allocates from a declared count
// without validating that count against bytes actually received) and depth-bounded. This rule applies
// per direction: encoding a response is not gated by it, because the encoder's input is a service-owned
// domain object rather than caller-controlled bytes. Request/response size limits bound input bytes
// only, so they do not mitigate amplification and are not a substitute for these bounds.
//
// json, yaml and protobuf satisfy this rule. msgpack and gob do not, so [Media.CanDecodeRequest] and
// [Content.NewStreamFromContentType] reject them for request decoding even though both remain valid
// media types and valid response codecs (see [Content.NewFromMedia] and the go-service HTTP client,
// which decodes responses through the same registries without this restriction).
//
// Start with [NewContent], [Content.NewFromRequest], and [Content.NewFromMedia].
package content
