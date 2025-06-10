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

// RegisterParams for rpc.
type RegisterParams struct {
	di.In

	Mux     *http.ServeMux
	Content *content.Content
	Pool    *sync.BufferPool
}

// Register for rpc.
func Register(params RegisterParams) {
	mux = params.Mux
	cont = params.Content
	pool = params.Pool
}
