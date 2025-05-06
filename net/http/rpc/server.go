package rpc

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/strings"
)

// Route for rpc.
func Route[Req any, Res any](path string, handler content.RequestHandler[Req, Res]) {
	mux.HandleFunc(strings.Join(" ", http.MethodPost, path), content.NewRequestHandler(cont, "rpc", handler))
}
