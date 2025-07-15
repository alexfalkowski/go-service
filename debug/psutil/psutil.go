package psutil

import (
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
)

// Register for debug.
func Register(name env.Name, cont *content.Content, mux *http.ServeMux) {
	mux.HandleFunc(http.Pattern(name, "/debug/psutil"), NewHandler(cont))
}
