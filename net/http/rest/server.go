package rest

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewServer constructs a Server that registers REST routes on router using unaryContent and streamContent
// for encoding/decoding, and opts for streaming route behavior.
//
// NewServer is typically called via dependency injection during application startup.
func NewServer(router *http.Router, unaryContent *unary.Content, streamContent *stream.Content, opts stream.Options) *Server {
	return &Server{router: router, unaryContent: unaryContent, streamContent: streamContent, opts: opts}
}

// Server registers REST-style HTTP handlers on a router.
//
// The streaming options are passed unchanged to this type's streaming route helpers (StreamRoute, StreamGet,
// StreamRouteRequest, StreamPost, StreamPut, StreamPatch) via
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.NewHandler] and
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.NewRequestHandler]; see those
// constructors for the timeout and per-value receive-size semantics.
type Server struct {
	router        *http.Router
	unaryContent  *unary.Content
	streamContent *stream.Content
	opts          stream.Options
}

// Delete registers an HTTP DELETE handler under pattern.
//
// The effective route pattern passed to the router is a method-qualified pattern of the form:
//
//	"<METHOD> <pattern>"
//
// For example:
//
//	server.Delete("/health", handler) // registers "DELETE /health"
//
// This method delegates to Route.
func (s *Server) Delete[Res any](pattern string, handler unary.Handler[Res], opts ...http.RouteOption) {
	s.Route(strings.Join(strings.Space, http.MethodDelete, pattern), handler, opts...)
}

// Get registers an HTTP GET handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This method delegates to Route.
func (s *Server) Get[Res any](pattern string, handler unary.Handler[Res], opts ...http.RouteOption) {
	s.Route(strings.Join(strings.Space, http.MethodGet, pattern), handler, opts...)
}

// Post registers an HTTP POST handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This method delegates to RouteRequest.
func (s *Server) Post[Req any, Res any](pattern string, handler unary.RequestHandler[Req, Res], opts ...http.RouteOption) {
	s.RouteRequest(strings.Join(strings.Space, http.MethodPost, pattern), handler, opts...)
}

// Put registers an HTTP PUT handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This method delegates to RouteRequest.
func (s *Server) Put[Req any, Res any](pattern string, handler unary.RequestHandler[Req, Res], opts ...http.RouteOption) {
	s.RouteRequest(strings.Join(strings.Space, http.MethodPut, pattern), handler, opts...)
}

// Patch registers an HTTP PATCH handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This method delegates to RouteRequest.
func (s *Server) Patch[Req any, Res any](pattern string, handler unary.RequestHandler[Req, Res], opts ...http.RouteOption) {
	s.RouteRequest(strings.Join(strings.Space, http.MethodPatch, pattern), handler, opts...)
}

// RouteRequest registers a handler under pattern that decodes a request and encodes a response.
//
// The handler is built using [github.com/alexfalkowski/go-service/v2/net/http/content/unary.NewRequestHandler], which:
//   - decodes the request body using Content-Type, falling back to JSON when Content-Type is absent and
//     rejecting it with 415 when it is unparseable, unregistered, or intentionally undecodable,
//   - decodes the request body into a newly allocated request model, and
//   - encodes the returned response model using the negotiated media type.
//
// Registration:
// The resulting handler is registered on the configured router.
// Options are forwarded to the router registration.
func (s *Server) RouteRequest[Req any, Res any](pattern string, handler unary.RequestHandler[Req, Res], opts ...http.RouteOption) {
	s.router.HandleRoute(pattern, unary.NewRequestHandler(s.unaryContent, handler), opts...)
}

// Route registers a handler under pattern that encodes a response.
//
// The handler is built using [github.com/alexfalkowski/go-service/v2/net/http/content/unary.NewHandler], which:
//   - selects an encoder based on the first Accept media type, falling back to Content-Type when
//     Accept is absent, and
//   - encodes the returned response model using the negotiated media type.
//
// Registration:
// The resulting handler is registered on the configured router.
// Options are forwarded to the router registration.
func (s *Server) Route[Res any](pattern string, handler unary.Handler[Res], opts ...http.RouteOption) {
	s.router.HandleRoute(pattern, unary.NewHandler(s.unaryContent, handler), opts...)
}

// StreamGet registers an HTTP GET handler under pattern for a send-only streaming response: the
// response streams incrementally rather than being buffered whole.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This method delegates to StreamRoute.
func (s *Server) StreamGet[Res any](pattern string, handler stream.Handler[Res], opts ...http.RouteOption) {
	s.StreamRoute(strings.Join(strings.Space, http.MethodGet, pattern), handler, opts...)
}

// StreamPost registers an HTTP POST handler under pattern for a bidirectional stream: both the
// request and response bodies stream incrementally rather than being buffered whole.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This method delegates to StreamRouteRequest.
//
// HTTP/2 requirement: see StreamRouteRequest. StreamGet and StreamRoute have no such requirement and
// stay fully supported on HTTP/1.1.
func (s *Server) StreamPost[Req any, Res any](pattern string, handler stream.RequestHandler[Req, Res], opts ...http.RouteOption) {
	s.StreamRouteRequest(strings.Join(strings.Space, http.MethodPost, pattern), handler, opts...)
}

// StreamPut registers an HTTP PUT handler under pattern for a bidirectional stream: both the request
// and response bodies stream incrementally rather than being buffered whole.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This method delegates to StreamRouteRequest.
//
// HTTP/2 requirement: see StreamRouteRequest.
func (s *Server) StreamPut[Req any, Res any](pattern string, handler stream.RequestHandler[Req, Res], opts ...http.RouteOption) {
	s.StreamRouteRequest(strings.Join(strings.Space, http.MethodPut, pattern), handler, opts...)
}

// StreamPatch registers an HTTP PATCH handler under pattern for a bidirectional stream: both the
// request and response bodies stream incrementally rather than being buffered whole.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This method delegates to StreamRouteRequest.
//
// HTTP/2 requirement: see StreamRouteRequest.
func (s *Server) StreamPatch[Req any, Res any](pattern string, handler stream.RequestHandler[Req, Res], opts ...http.RouteOption) {
	s.StreamRouteRequest(strings.Join(strings.Space, http.MethodPatch, pattern), handler, opts...)
}

// StreamRouteRequest registers a handler under pattern for a bidirectional stream: both the request
// and response bodies stream incrementally rather than being buffered whole.
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
func (s *Server) StreamRouteRequest[Req any, Res any](pattern string, handler stream.RequestHandler[Req, Res], opts ...http.RouteOption) {
	s.router.HandleRoute(pattern, stream.NewRequestHandler(s.streamContent, s.opts, handler), append(opts, http.WithRouteStreaming())...)
}

// StreamRoute registers a handler under pattern for a send-only streaming response: the response
// streams incrementally rather than being buffered whole; the request body is not streamed.
//
// The handler is built using [github.com/alexfalkowski/go-service/v2/net/http/content/stream.NewHandler], which:
//   - selects a streaming encoder based on Accept, falling back to Content-Type when Accept is absent,
//     via the streaming media registry, and
//   - rejects an Accept that cannot be satisfied with 406.
//
// Unlike StreamRouteRequest and its method-qualified helpers, StreamRoute has no HTTP/2 requirement
// and stays fully supported on HTTP/1.1 chunked responses.
//
// Registration:
// The resulting handler is registered on the configured router, and the route is marked streaming on
// the router's route policy (see
// [github.com/alexfalkowski/go-service/v2/net/http.WithRouteResponseStreaming]). Its request body
// retains the usual cumulative request-body limit because it is not streamed.
// Options are additive to the route's response streaming policy.
func (s *Server) StreamRoute[Res any](pattern string, handler stream.Handler[Res], opts ...http.RouteOption) {
	s.router.HandleRoute(pattern, stream.NewHandler(s.streamContent, s.opts, handler), append(opts, http.WithRouteResponseStreaming())...)
}
