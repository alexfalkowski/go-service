package rpc

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Route registers an RPC-style HTTP POST handler under pattern.
//
// The effective route pattern passed to the router is method-qualified and has the form:
//
//	"<METHOD> <pattern>"
//
// For example:
//
//	Route("/greet.v1.Greeter/SayHello", handler) // registers "POST /greet.v1.Greeter/SayHello"
//
// The handler is constructed using [github.com/alexfalkowski/go-service/v2/net/http/content.NewRequestHandler], which:
//   - decodes the request body using Content-Type, falling back to JSON when Content-Type is absent and
//     rejecting it with 415 when it is unparseable, unregistered, or intentionally undecodable,
//   - decodes the request body into a newly allocated request model, and
//   - encodes the returned response model using Accept, falling back to Content-Type when Accept is absent.
//
// Registration:
// The resulting handler is registered on the package-level router configured via [Register].
// [Register] must be called before Route; otherwise router/cont will be nil and this function will panic.
func Route[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	router.Handle(
		strings.Join(strings.Space, http.MethodPost, pattern),
		content.NewRequestHandler(cont, handler),
	)
}

// StreamRoute registers an RPC-style HTTP POST handler under pattern for a bidirectional stream: both
// the request and response bodies stream incrementally rather than being buffered whole.
//
// The effective route pattern passed to the router is method-qualified and has the form:
//
//	"POST <pattern>"
//
// Unlike REST, where a streamed response with no streamed request (StreamGet/StreamRoute) is a
// distinct shape, RPC's single POST-only style makes StreamRoute bidi-capable by nature, mirroring
// how [Route] is already POST-only.
//
// The handler is built using
// [github.com/alexfalkowski/go-service/v2/net/http/content.NewRequestStreamHandler], which:
//   - resolves the request decoder from Content-Type and the response encoder from the first Accept
//     media type (falling back to Content-Type), both via the streaming media registry, and
//   - rejects an unregistered or unparseable streaming media type with 415.
//
// HTTP/2 requirement:
// Bidirectional streaming requires HTTP/2 (including h2c): an HTTP/1.x request body is buffered ahead
// of the handler by intermediaries and the Go transport, so interleaving Recv/Send hangs rather than
// failing. [content.NewRequestStreamHandler] rejects requests with req.ProtoMajor < 2 with
// 505 HTTP Version Not Supported before the handler runs. Deploying a route registered through this
// helper therefore requires the server, and any intermediary in front of it, to support HTTP/2 or h2c.
//
// Registration:
// The resulting handler is registered on the package-level router configured via [Register], and the
// route is marked streaming on the router's route policy (see
// [github.com/alexfalkowski/go-service/v2/net/http.RoutePolicy.Streaming]) so inbound request body
// limiting is applied lazily instead of buffering the whole body.
// [Register] must be called before StreamRoute; otherwise router/cont will be nil and this function
// will panic.
//
// Inbound size limiting:
// maxReceiveSize (set via [Register]) bounds each value decoded by the resulting stream's
// Recv, not the request body as a whole — see [content.NewRequestStreamHandler].
func StreamRoute[Req any, Res any](pattern string, handler content.RequestStreamHandler[Req, Res]) {
	router.HandleStreaming(
		strings.Join(strings.Space, http.MethodPost, pattern),
		content.NewRequestStreamHandler(cont, timeout, maxReceiveSize, handler),
	)
}
