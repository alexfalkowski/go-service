package grpc

import (
	"github.com/alexfalkowski/go-service/v2/os"
)

var fs *os.FS

// Register injects the filesystem dependency used by this package.
//
// The registered filesystem is consulted when constructing TLS configuration for both clients and servers.
// It is used to resolve certificate/key "source strings" (for example `file:/path/to/cert` or `env:VAR`) via
// the `os.FS.ReadSource` helper when materializing `*tls.Config`.
//
// Registration is required because the filesystem is stored in a package-level variable. If TLS is enabled,
// call Register during application startup (composition root) before constructing any gRPC clients or servers;
// otherwise TLS config construction may fail to load key material.
func Register(f *os.FS) {
	fs = f
}
