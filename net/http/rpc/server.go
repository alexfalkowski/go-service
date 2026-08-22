package rpc

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewServer constructs a Server that registers RPC routes on router using unaryContent and streamContent
// for encoding/decoding, and opts for streaming route behavior.
//
// NewServer is typically called via dependency injection during application startup.
func NewServer(router *http.Router, unaryContent *unary.Content, streamContent *stream.Content, opts stream.Options) *Server {
	return &Server{router: router, unaryContent: unaryContent, streamContent: streamContent, opts: opts}
}

// Server registers RPC-style HTTP handlers on a router.
//
// The streaming options are passed unchanged to StreamRoute via
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.NewRequestHandler]; see that
// constructor for the timeout and per-value receive-size semantics.
type Server struct {
	router        *http.Router
	unaryContent  *unary.Content
	streamContent *stream.Content
	opts          stream.Options
}

// Route registers an RPC-style HTTP POST handler under pattern.
//
// The effective route pattern passed to the router is method-qualified and has the form:
//
//	"<METHOD> <pattern>"
//
// For example:
//
//	server.Route("/greet.v1.Greeter/SayHello", handler) // registers "POST /greet.v1.Greeter/SayHello"
//
// The handler is constructed using [github.com/alexfalkowski/go-service/v2/net/http/content/unary.NewRequestHandler], which:
//   - decodes the request body using Content-Type, falling back to JSON when Content-Type is absent and
//     rejecting it with 415 when it is unparseable, unregistered, or intentionally undecodable,
//   - decodes the request body into a newly allocated request model, and
//   - encodes the returned response model using Accept, falling back to Content-Type when Accept is absent.
//
// Registration:
// The resulting handler is registered on the configured router.
// Options are forwarded to the router registration.
func (s *Server) Route[Req any, Res any](pattern string, handler unary.RequestHandler[Req, Res], opts ...http.RouteOption) {
	s.router.HandleRoute(
		strings.Join(strings.Space, http.MethodPost, pattern),
		unary.NewRequestHandler(s.unaryContent, handler),
		opts...,
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
// how [Server.Route] is already POST-only.
//
// The handler is built using
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.NewRequestHandler], which:
//   - resolves the request decoder from Content-Type, rejecting an unregistered or unparseable
//     streaming media type with 415, and
//   - resolves the response encoder from Accept (falling back to Content-Type), rejecting an Accept
//     that cannot be satisfied with 406.
//
// HTTP/2 requirement:
// Bidirectional streaming requires HTTP/2 (including h2c): an HTTP/1.x request body is buffered ahead
// of the handler by intermediaries and the Go transport, so interleaving Recv/Send hangs rather than
// failing. [stream.NewRequestHandler] rejects requests with req.ProtoMajor < 2 with
// 505 HTTP Version Not Supported before the handler runs. Deploying a route registered through this
// method therefore requires the server, and any intermediary in front of it, to support HTTP/2 or h2c.
//
// Registration:
// The resulting handler is registered on the configured router, and the route is marked streaming on
// the router's route policy (see
// [github.com/alexfalkowski/go-service/v2/net/http.WithRouteStreaming]) so inbound request body
// limiting is applied lazily instead of buffering the whole body.
// Options are additive to the route's bidirectional streaming policy.
//
// Inbound size limiting:
// The configured maximum receive size bounds each value decoded by the resulting stream's Recv, not the
// request body as a whole — see [stream.NewRequestHandler].
func (s *Server) StreamRoute[Req any, Res any](pattern string, handler stream.RequestHandler[Req, Res], opts ...http.RouteOption) {
	s.router.HandleRoute(
		strings.Join(strings.Space, http.MethodPost, pattern),
		stream.NewRequestHandler(s.streamContent, s.opts, handler),
		append(opts, http.WithRouteStreaming())...,
	)
}
