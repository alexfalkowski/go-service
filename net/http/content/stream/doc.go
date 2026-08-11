// Package stream provides NDJSON HTTP request and response streaming on top of content codecs.
//
// [Content] owns a streaming codec registry that is separate from the single-value registry in
// [github.com/alexfalkowski/go-service/v2/net/http/content/unary]. Streaming media negotiation is strict:
// an unregistered or unparseable declared media type is rejected instead of falling back to another
// representation. With no Content-Type or Accept header, response negotiation selects NDJSON.
//
// # Handlers
//
// [NewHandler] creates a send-only streaming response handler and works over HTTP/1.x or HTTP/2.
// [NewRequestHandler] creates a bidirectional request/response handler and requires HTTP/2 (including
// h2c). Use one goroutine at a time with [Stream.Send]; for [RequestStream], the supported pattern is
// one goroutine calling Send and one calling Recv. See the executable streaming examples in
// `net/http/rest` and `net/http/client` for complete server and client calls.
//
// # Limits and lifecycle
//
// [Options.MaxReceiveSize] caps each decoded request value, not the cumulative request stream. Request
// decoding accepts only codecs that are safe for untrusted input, matching the decoder-bounds rule in
// the unary package documentation. [Options.ReadTimeout] and [Options.WriteTimeout] are per-message
// inactivity budgets, while [Options.Drain] cancels handlers during server shutdown. Handlers waiting
// on an upstream source must return after their context is canceled.
package stream
