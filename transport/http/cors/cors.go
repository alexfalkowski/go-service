package cors

import (
	"github.com/rs/cors"
)

// New cors.
func New() *cors.Cors {
	return cors.AllowAll()
}
