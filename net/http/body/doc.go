// Package body provides request-body handling helpers for HTTP handlers.
//
// [NewHandler] buffers the whole body and enforces a cumulative size limit before the handler runs.
// [NewLazyHandler] never buffers and does not enforce that limit as the body is read; it only keeps
// NewHandler's declared-Content-Length short-circuit, for routes that read their body incrementally as
// it arrives (streaming routes) and cap individual decoded values themselves instead of a cumulative
// total.
package body
