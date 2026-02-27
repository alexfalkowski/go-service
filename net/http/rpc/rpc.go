package rpc

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/sync"
)

var (
	mux  *http.ServeMux
	cont *content.Content
	pool *sync.BufferPool
)

// RegisterParams defines dependencies used to register RPC package globals.
//
// This package relies on package-level registration to avoid threading commonly shared dependencies
// through every helper call. These dependencies are typically provided by DI wiring.
type RegisterParams struct {
	di.In

	// Mux is the HTTP mux where server-side routes will be registered by this package's helpers.
	Mux *http.ServeMux

	// Content resolves encoders/decoders based on HTTP media types (Content-Type).
	Content *content.Content

	// Pool is a buffer pool used by client helpers to reduce allocations while encoding/decoding bodies.
	Pool *sync.BufferPool
}

// Register stores the dependencies used by server and client helpers in package-level variables.
//
// Register is expected to be called during application startup (typically via dependency injection).
//
// Important: Register must be called before using any server-side route helpers or client helpers in this
// package. If it is not called, globals will be nil and helper calls will panic.
func Register(params RegisterParams) {
	mux = params.Mux
	cont = params.Content
	pool = params.Pool
}
