package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/strings"
)

// Path will strip / from the start.
func Path(req *http.Request) string {
	path := req.URL.Path
	if strings.IsEmpty(path) {
		return path
	}

	return path[1:]
}
