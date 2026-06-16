package transport

import (
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
)

// Config configures server-side transport wiring for a service.
//
// It is a top-level configuration object that enables and configures the supported server
// transport stacks (currently HTTP and gRPC). Each nested transport config is optional:
// when nil or disabled, the corresponding server constructors typically return nil and no
// server is created.
//
// Client construction uses explicit client options and client-side config types such as
// [github.com/alexfalkowski/go-service/v2/config/client.Config]; transport.http and
// transport.grpc are the server-side config trees.
//
// The struct tags are compatible with the repository's config decoder (YAML/JSON/TOML).
type Config struct {
	// Access configures shared authorization policy for all enabled transport server stacks.
	//
	// When nil, transport access control is disabled. When configured, the same controller is injected into
	// HTTP and gRPC server middleware/interceptors, and policies distinguish protocols through
	// meta.TransportServiceMethod values such as "http:GET /users/{id}" or
	// "grpc:/package.Service/Method".
	Access *access.Config `yaml:"access,omitempty" json:"access,omitempty" toml:"access,omitempty"`

	// GRPC configures the gRPC server transport stack.
	//
	// When nil or when the nested config is disabled, gRPC transport wiring is effectively
	// turned off and constructors such as [github.com/alexfalkowski/go-service/v2/transport/grpc.NewServer] typically return nil.
	GRPC *grpc.Config `yaml:"grpc,omitempty" json:"grpc,omitempty" toml:"grpc,omitempty"`

	// HTTP configures the HTTP server transport stack.
	//
	// When nil or when the nested config is disabled, HTTP transport wiring is effectively
	// turned off and constructors such as [github.com/alexfalkowski/go-service/v2/transport/http.NewServer] typically return nil.
	HTTP *http.Config `yaml:"http,omitempty" json:"http,omitempty" toml:"http,omitempty"`
}

// IsEnabled reports whether transport configuration is present.
//
// This is a nil-safe convenience used by callers that treat an omitted transport config as "disabled".
// Note that this method only checks that the top-level transport config is non-nil; individual transports
// (HTTP/gRPC) are enabled/disabled by their own nested configuration.
func (c *Config) IsEnabled() bool {
	return c != nil
}
