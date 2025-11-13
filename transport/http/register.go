package http

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/os"
)

var (
	name env.Name
	fs   *os.FS
)

// Register for http.
func Register(n env.Name, f *os.FS) {
	name = n
	fs = f
}
