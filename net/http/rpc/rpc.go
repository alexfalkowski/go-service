package rpc

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-sync"
)

var (
	router        *http.Router
	unaryContent  *unary.Content
	pool          *sync.BufferPool
	streamContent *stream.Content
	opts          stream.Options
)

// Register stores the dependencies used by server and client helpers in package-level variables.
//
// Register is expected to be called during application startup (typically via dependency injection).
//
// Important: Register must be called before using any server-side route helpers or client helpers in this
// package. If it is not called, globals will be nil and helper calls will panic.
//
// unaryContent resolves single-value media types, and streamContent resolves streaming media types for StreamRoute
// and the embedded HTTP client. p also buffers client response bodies.
//
// so is passed through unchanged to this package's streaming route helper (StreamRoute) via
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.NewRequestHandler]; see that
// constructor for the timeout and per-value receive-size semantics.
func Register(r *http.Router, uc *unary.Content, sc *stream.Content, p *sync.BufferPool, so stream.Options) {
	router = r
	unaryContent = uc
	pool = p
	streamContent = sc
	opts = so
}
