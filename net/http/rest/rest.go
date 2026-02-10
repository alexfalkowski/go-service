package rest

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

// RegisterParams defines dependencies used to register REST globals.
type RegisterParams struct {
	di.In
	Mux     *http.ServeMux
	Content *content.Content
	Pool    *sync.BufferPool
}

// Register stores the dependencies used by server and client helpers.
//
// Register is expected to be called during application startup (typically via Fx). Server-side route helpers use the
// registered mux, and client helpers use the registered content codecs and buffer pool.
func Register(params RegisterParams) {
	mux = params.Mux
	cont = params.Content
	pool = params.Pool
}
