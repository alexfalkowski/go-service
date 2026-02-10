package grpc

import (
	"github.com/alexfalkowski/go-service/v2/os"
)

var fs *os.FS

// Register injects the filesystem dependency used by this package.
//
// The registered filesystem is used when constructing TLS configuration (for example to read certificates/keys)
// in client and server constructors. If TLS is enabled, ensure Register is called during application startup
// before constructing clients/servers.
func Register(f *os.FS) {
	fs = f
}
