// Package body provides HTTP request body size-limit middleware.
//
// The middleware caps inbound request bodies before downstream handlers decode
// them. Accepted requests receive a buffered replacement body, while oversized
// or unreadable bodies are rejected through the response writer. [NewHandler] never
// buffers routes marked streaming in its route policy: it keeps the
// declared-Content-Length short-circuit but does not enforce a cumulative limit as the body is read,
// since a streaming route caps individual decoded values at the content layer instead of a whole-body
// total.
//
// Start with [NewHandler].
package body
