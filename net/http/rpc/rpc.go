package rpc

import (
	encodingstream "github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	contentstream "github.com/alexfalkowski/go-service/v2/net/http/content/stream"
	"github.com/alexfalkowski/go-sync"
)

var (
	router *http.Router
	cont   *content.Content
	sm     *encodingstream.Map
	pool   *sync.BufferPool
	opts   contentstream.Options
)

// Register stores the dependencies used by server and client helpers in package-level variables.
//
// Register is expected to be called during application startup (typically via dependency injection).
//
// Important: Register must be called before using any server-side route helpers or client helpers in this
// package. If it is not called, globals will be nil and helper calls will panic.
//
// c resolves single-value media types, while s resolves streaming media types for StreamRoute and
// the embedded HTTP client.
//
// so is passed through unchanged to this package's streaming route helper (StreamRoute) via
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.NewRequestHandler]; see that
// constructor for the timeout and per-value receive-size semantics.
func Register(r *http.Router, c *content.Content, s *encodingstream.Map, p *sync.BufferPool, so contentstream.Options) {
	router = r
	cont = c
	sm = s
	pool = p
	opts = so
}
