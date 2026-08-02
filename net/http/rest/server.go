package rest

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Delete registers an HTTP DELETE handler under pattern.
//
// The effective route pattern passed to the router is a method-qualified pattern of the form:
//
//	"<METHOD> <pattern>"
//
// For example:
//
//	Delete("/health", handler) // registers "DELETE /health"
//
// This helper delegates to Route.
func Delete[Res any](pattern string, handler unary.Handler[Res]) {
	Route(strings.Join(strings.Space, http.MethodDelete, pattern), handler)
}

// Get registers an HTTP GET handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to Route.
func Get[Res any](pattern string, handler unary.Handler[Res]) {
	Route(strings.Join(strings.Space, http.MethodGet, pattern), handler)
}

// Post registers an HTTP POST handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to RouteRequest.
func Post[Req any, Res any](pattern string, handler unary.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(strings.Space, http.MethodPost, pattern), handler)
}

// Put registers an HTTP PUT handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to RouteRequest.
func Put[Req any, Res any](pattern string, handler unary.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(strings.Space, http.MethodPut, pattern), handler)
}

// Patch registers an HTTP PATCH handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to RouteRequest.
func Patch[Req any, Res any](pattern string, handler unary.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(strings.Space, http.MethodPatch, pattern), handler)
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
// The resulting handler is registered on the package-level router configured via [Register].
// [Register] must be called before RouteRequest; otherwise router/unaryContent will be nil and this function will panic.
func RouteRequest[Req any, Res any](pattern string, handler unary.RequestHandler[Req, Res]) {
	router.Handle(pattern, unary.NewRequestHandler(unaryContent, handler))
}

// Route registers a handler under pattern that encodes a response.
//
// The handler is built using [github.com/alexfalkowski/go-service/v2/net/http/content/unary.NewHandler], which:
//   - selects an encoder based on the first Accept media type, falling back to Content-Type when
//     Accept is absent, and
//   - encodes the returned response model using the negotiated media type.
//
// Registration:
// The resulting handler is registered on the package-level router configured via Register.
// Register must be called before Route; otherwise router/unaryContent will be nil and this function will panic.
func Route[Res any](pattern string, handler unary.Handler[Res]) {
	router.Handle(pattern, unary.NewHandler(unaryContent, handler))
}

// StreamGet registers an HTTP GET handler under pattern for a send-only streaming response: the
// response streams incrementally rather than being buffered whole.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to StreamRoute.
func StreamGet[Res any](pattern string, handler stream.Handler[Res]) {
	StreamRoute(strings.Join(strings.Space, http.MethodGet, pattern), handler)
}

// StreamPost registers an HTTP POST handler under pattern for a bidirectional stream: both the
// request and response bodies stream incrementally rather than being buffered whole.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to StreamRouteRequest.
//
// HTTP/2 requirement: see StreamRouteRequest. StreamGet and StreamRoute have no such requirement and
// stay fully supported on HTTP/1.1.
func StreamPost[Req any, Res any](pattern string, handler stream.RequestHandler[Req, Res]) {
	StreamRouteRequest(strings.Join(strings.Space, http.MethodPost, pattern), handler)
}

// StreamPut registers an HTTP PUT handler under pattern for a bidirectional stream: both the request
// and response bodies stream incrementally rather than being buffered whole.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to StreamRouteRequest.
//
// HTTP/2 requirement: see StreamRouteRequest.
func StreamPut[Req any, Res any](pattern string, handler stream.RequestHandler[Req, Res]) {
	StreamRouteRequest(strings.Join(strings.Space, http.MethodPut, pattern), handler)
}

// StreamPatch registers an HTTP PATCH handler under pattern for a bidirectional stream: both the
// request and response bodies stream incrementally rather than being buffered whole.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to StreamRouteRequest.
//
// HTTP/2 requirement: see StreamRouteRequest.
func StreamPatch[Req any, Res any](pattern string, handler stream.RequestHandler[Req, Res]) {
	StreamRouteRequest(strings.Join(strings.Space, http.MethodPatch, pattern), handler)
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
// helper therefore requires the server, and any intermediary in front of it, to support HTTP/2 or h2c.
//
// Registration:
// The resulting handler is registered on the package-level router configured via [Register], and the
// route is marked streaming on the router's route policy (see
// [github.com/alexfalkowski/go-service/v2/net/http.RoutePolicy.Streaming]) so inbound request body
// limiting is applied lazily instead of buffering the whole body.
// [Register] must be called before StreamRouteRequest; otherwise router/streamContent will be nil and this
// function will panic.
//
// Inbound size limiting:
// opts.MaxReceiveSize (set via [Register]) bounds each value decoded by the resulting stream's
// Recv, not the request body as a whole — see [stream.NewRequestHandler].
func StreamRouteRequest[Req any, Res any](pattern string, handler stream.RequestHandler[Req, Res]) {
	router.HandleStreaming(pattern, stream.NewRequestHandler(streamContent, opts, handler))
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
// The resulting handler is registered on the package-level router configured via Register, and the
// route is marked streaming on the router's route policy (see
// [github.com/alexfalkowski/go-service/v2/net/http.RoutePolicy.Streaming]) so inbound request body
// limiting is applied lazily instead of buffering the whole body.
// Register must be called before StreamRoute; otherwise router/streamContent will be nil and this function will
// panic.
func StreamRoute[Res any](pattern string, handler stream.Handler[Res]) {
	router.HandleStreaming(pattern, stream.NewHandler(streamContent, opts, handler))
}
