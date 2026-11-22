package http

import (
	"github.com/alexfalkowski/go-service/v3/os"
)

var fs *os.FS

// Register injects the filesystem dependency used by this package.
//
// The registered filesystem is consulted when constructing TLS configuration for both clients and servers.
// It is used by [github.com/alexfalkowski/go-service/v3/config/server.NewConfig] and [github.com/alexfalkowski/go-service/v3/config/client.NewConfig] to resolve TLS
// "source strings" (for example `file:<path>` or `env:VAR`) via the
// [os.FS.ReadSource] helper when materializing a runtime `*crypto/tls.Config`.
//
// The standard [github.com/alexfalkowski/go-service/v3/module.Server] and
// [github.com/alexfalkowski/go-service/v3/module.Client] bundles register the filesystem before constructing OpenFeature
// providers. [Module] also registers it for direct HTTP transport composition. For manual setup, call
// Register before constructing a TLS client or server that reads file sources.
func Register(f *os.FS) {
	fs = f
}
