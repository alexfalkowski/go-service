package rest

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/sync"
	"go.uber.org/fx"
)

var (
	mux  *http.ServeMux
	cont *content.Content
	pool *sync.BufferPool
)

// RegisterParams for rest.
type RegisterParams struct {
	fx.In

	Mux     *http.ServeMux
	Content *content.Content
	Pool    *sync.BufferPool
}

// Register for rest.
func Register(params RegisterParams) {
	mux = params.Mux
	cont = params.Content
	pool = params.Pool
}
