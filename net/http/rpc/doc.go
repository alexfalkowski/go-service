// Package rpc provides RPC-style HTTP handler registration and client helpers for go-service.
//
// This package is built on top of net/http/content. It relies on package-level registration (see Register)
// to supply the HTTP mux, content codec helpers, and buffer pool that are used when wiring handlers and clients.
//
// Server-side helpers register POST handlers on the configured mux.
// Client helpers (NewClient) build a net/http/client.Client using the registered content and buffer pool.
//
// Note: Register must be called before using server or client helpers; otherwise globals will be nil and
// handler/client construction will panic.
package rpc
