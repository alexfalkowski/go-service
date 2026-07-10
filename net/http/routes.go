package http

import "github.com/alexfalkowski/go-service/v2/strings"

// NewRoutePolicy constructs an empty route policy registry.
func NewRoutePolicy() *RoutePolicy {
	return &RoutePolicy{
		operations:      map[string]struct{}{},
		unauthenticated: map[string]struct{}{},
	}
}

// RoutePolicy stores route policy used by HTTP transport middleware.
//
// Route registration code records policy here so middleware can make exact route decisions without
// inferring intent from path substrings. RoutePolicy is intended to be populated during startup before serving
// requests.
type RoutePolicy struct {
	operations      map[string]struct{}
	unauthenticated map[string]struct{}
}

// Operation marks pattern as a service-owned operation path.
//
// Operation matching is path-only so method mismatches can reach the mux and receive normal method handling.
func (r *RoutePolicy) Operation(pattern string) {
	r.operations[routePatternPath(pattern)] = struct{}{}
}

// AllowUnauthenticated marks pattern as not requiring transport token authentication.
func (r *RoutePolicy) AllowUnauthenticated(pattern string) {
	r.unauthenticated[pattern] = struct{}{}
}

// IsOperation reports whether req targets a registered operation path.
func (r *RoutePolicy) IsOperation(req *Request) bool {
	_, ok := r.operations[req.URL.Path]
	return ok
}

// IsUnauthenticated reports whether req targets a route that does not require transport token authentication.
func (r *RoutePolicy) IsUnauthenticated(req *Request) bool {
	if !strings.IsEmpty(req.Pattern) {
		_, ok := r.unauthenticated[req.Pattern]
		return ok
	}

	if _, ok := r.unauthenticated[routeRequestPattern(req)]; ok {
		return true
	}

	_, ok := r.unauthenticated[req.URL.Path]
	return ok
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

// Handle registers handler for pattern on the Router's mux.
func (r *Router) Handle(pattern string, handler Handler) {
	r.mux.Handle(pattern, handler)
}

// HandleOperation registers handler and marks pattern as a service-owned operation path.
func (r *Router) HandleOperation(pattern string, handler Handler) {
	r.routePolicy.Operation(pattern)
	r.Handle(pattern, handler)
}

// HandleOperationFunc registers handler and marks pattern as a service-owned operation path.
func (r *Router) HandleOperationFunc(pattern string, handler HandlerFunc) {
	r.HandleOperation(pattern, handler)
}

// HandleUnauthenticated registers handler and marks pattern as not requiring transport token authentication.
func (r *Router) HandleUnauthenticated(pattern string, handler Handler) {
	r.routePolicy.AllowUnauthenticated(pattern)
	r.Handle(pattern, handler)
}

func routePatternPath(pattern string) string {
	return strings.CutAfter(pattern, strings.Space)
}

func routeRequestPattern(req *Request) string {
	return strings.Join(strings.Space, req.Method, req.URL.Path)
}
