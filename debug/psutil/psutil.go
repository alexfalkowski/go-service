package psutil

import (
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
)

// Register for debug.
func Register(mux *http.ServeMux, cont *content.Content) {
	mux.HandleFunc("/debug/psutil", NewHandler(cont))
}
