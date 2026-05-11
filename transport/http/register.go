package http

import (
	"github.com/alexfalkowski/go-service/v2/os"
)

var fs *os.FS

// Register injects the filesystem dependency used by this package.
//
// The registered filesystem is consulted when constructing TLS configuration for both clients and servers.
// It is used by `config/server.NewConfig` and `config/client.NewConfig` to resolve TLS
// "source strings" (for example `file:/path/to/cert` or `env:VAR`) via the
// `os.FS.ReadSource` helper when materializing a runtime `*crypto/tls.Config`.
//
// In the standard go-service module graph, this registration is performed automatically by `Module`
// (and therefore by higher-level bundles such as `module.Server`). Call Register directly only when you
// intentionally compose HTTP transport pieces manually outside that graph.
//
// If TLS is enabled and registration has not happened before constructing HTTP clients or servers, TLS
// config construction may fail to load key material.
func Register(f *os.FS) {
	fs = f
}
