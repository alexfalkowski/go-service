package rpc

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Route for rpc.
func Route[Req any, Res any](path string, handler content.RequestHandler[Req, Res]) {
	mux.HandleFunc(strings.Join(" ", http.MethodPost, path), content.NewRequestHandler(cont, handler))
}
