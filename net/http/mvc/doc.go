// Package mvc provides a small MVC-style HTML rendering layer for go-service HTTP servers.
//
// This package provides helpers to register controller-driven routes that render HTML templates.
// It relies on package-level registration (see Register) to supply the HTTP mux, template function map,
// filesystem, and layout.
//
// Registration requirements:
//
// Call Register (typically via Fx) before using Route/Get/Post/etc. The routing helpers return false
// when MVC is not defined (for example when no filesystem or layout is registered).
//
// Context requirements:
//
// Controllers and views use net/http/meta to retrieve the request and response writer from context.
// The routing helpers set these values via meta.WithRequest and meta.WithResponse before invoking the controller.
package mvc
