package http

import "github.com/alexfalkowski/go-service/v2/strings"

// RouteOption configures an HTTP route.
type RouteOption func(*routeOptions)

type routeOptions struct {
	operation         bool
	unauthenticated   bool
	requestStreaming  bool
	responseStreaming bool
}

// WithRouteOperation marks a route as a service-owned operation path, for infrastructure endpoints such
// as the health and metrics routes registered by transport/http/health and transport/http/telemetry/metrics.
//
// Operation matching is path-only so method mismatches can reach the mux and receive normal method handling.
// Registration therefore marks the pattern's path for every method: registering "GET /admin" as an operation
// also marks a "POST /admin" registered by a separate HandleRoute call without this option.
//
// Supported middleware treats an operation route as exempt from transport token verification, access control,
// and rate limiting, and omits its per-request outcome log line (a recovered panic is still logged).
// WithRouteUnauthenticated instead bypasses transport token verification and access control while retaining
// rate limiting and normal outcome logging; use it only when another boundary protects the route.
func WithRouteOperation() RouteOption {
	return func(options *routeOptions) {
		options.operation = true
	}
}

// WithRouteUnauthenticated marks a route to bypass transport token verification and access control.
//
// It retains rate limiting and normal outcome logging, so use it only when another boundary protects the route.
func WithRouteUnauthenticated() RouteOption {
	return func(options *routeOptions) {
		options.unauthenticated = true
	}
}

// WithRouteRequestStreaming marks a route whose request body is streamed incrementally rather than buffered whole.
func WithRouteRequestStreaming() RouteOption {
	return func(options *routeOptions) {
		options.requestStreaming = true
	}
}

// WithRouteResponseStreaming marks a route whose response body is streamed incrementally rather than buffered whole.
func WithRouteResponseStreaming() RouteOption {
	return func(options *routeOptions) {
		options.responseStreaming = true
	}
}

// WithRouteStreaming marks a route whose request and response bodies are streamed incrementally rather than buffered whole.
func WithRouteStreaming() RouteOption {
	return func(options *routeOptions) {
		options.requestStreaming = true
		options.responseStreaming = true
	}
}

func options(opts ...RouteOption) *routeOptions {
	options := &routeOptions{}
	for _, opt := range opts {
		opt(options)
	}

	return options
}

// NewRoutePolicy constructs an empty route policy registry.
func NewRoutePolicy() *RoutePolicy {
	return &RoutePolicy{
		operations: map[string]struct{}{},
		routes:     map[string]routePolicy{},
	}
}

// RoutePolicy stores route policy used by HTTP transport middleware.
//
// Route registration code records policy here so middleware can make exact route decisions without
// inferring intent from path substrings. RoutePolicy is intended to be populated during startup before serving
// requests.
type RoutePolicy struct {
	operations map[string]struct{}
	routes     map[string]routePolicy
}

type routePolicy struct {
	unauthenticated   bool
	requestStreaming  bool
	responseStreaming bool
}

// IsOperation reports whether req targets a registered operation path.
func (r *RoutePolicy) IsOperation(req *Request) bool {
	if !strings.IsEmpty(req.Pattern) {
		if _, ok := r.operations[routePatternPath(req.Pattern)]; ok {
			return true
		}
	}

	_, ok := r.operations[req.URL.Path]
	return ok
}

// IsUnauthenticated reports whether req targets a route that does not require transport token authentication.
func (r *RoutePolicy) IsUnauthenticated(req *Request) bool {
	if !strings.IsEmpty(req.Pattern) {
		return r.routes[req.Pattern].unauthenticated
	}

	if r.routes[routeRequestPattern(req)].unauthenticated {
		return true
	}

	return r.routes[req.URL.Path].unauthenticated
}

// IsRequestStreaming reports whether req targets a route whose request body is streamed incrementally rather than buffered whole.
func (r *RoutePolicy) IsRequestStreaming(req *Request) bool {
	if !strings.IsEmpty(req.Pattern) {
		return r.routes[req.Pattern].requestStreaming
	}

	if r.routes[routeRequestPattern(req)].requestStreaming {
		return true
	}

	return r.routes[req.URL.Path].requestStreaming
}

// IsResponseStreaming reports whether req targets a route whose response body is streamed incrementally rather than buffered whole.
func (r *RoutePolicy) IsResponseStreaming(req *Request) bool {
	if !strings.IsEmpty(req.Pattern) {
		return r.routes[req.Pattern].responseStreaming
	}

	if r.routes[routeRequestPattern(req)].responseStreaming {
		return true
	}

	return r.routes[req.URL.Path].responseStreaming
}

// NewRouter constructs a Router backed by mux and routePolicy.
func NewRouter(mux *ServeMux, routePolicy *RoutePolicy) *Router {
	return &Router{mux: mux, routePolicy: routePolicy}
}

// Router registers HTTP handlers and their route policy on a mux.
type Router struct {
	mux         *ServeMux
	routePolicy *RoutePolicy
}

// HandleRoute applies options to the route policy and registers handler for pattern on the Router's mux.
func (r *Router) HandleRoute(pattern string, handler Handler, opts ...RouteOption) {
	options := options(opts...)

	if options.operation {
		r.routePolicy.operations[routePatternPath(pattern)] = struct{}{}
	}
	if options.unauthenticated || options.requestStreaming || options.responseStreaming {
		r.routePolicy.routes[pattern] = routePolicy{
			unauthenticated:   options.unauthenticated,
			requestStreaming:  options.requestStreaming,
			responseStreaming: options.responseStreaming,
		}
	}

	r.mux.Handle(pattern, handler)
}

func routePatternPath(pattern string) string {
	return strings.CutAfter(pattern, strings.Space)
}

func routeRequestPattern(req *Request) string {
	return strings.Join(strings.Space, req.Method, req.URL.Path)
}
