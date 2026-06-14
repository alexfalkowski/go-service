// Package body provides HTTP request body size-limit middleware.
//
// The middleware caps inbound request bodies before downstream handlers decode
// them. Accepted requests receive a buffered replacement body, while oversized
// or unreadable bodies are rejected through the response writer.
//
// Start with [NewHandler].
package body
