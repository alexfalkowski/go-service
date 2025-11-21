package http

import (
	"github.com/alexfalkowski/go-service/v2/os"
)

var fs *os.FS

// Register for http.
func Register(f *os.FS) {
	fs = f
}
