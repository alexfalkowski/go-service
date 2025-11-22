package rpc

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Route for rpc.
func Route[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	http.HandleFunc(
		mux,
		strings.Join(strings.Space, http.MethodPost, pattern),
		content.NewRequestHandler(cont, handler),
	)
}
