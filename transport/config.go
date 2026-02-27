package transport

import (
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
)

// Config configures the transport layer for a service.
//
// It is a top-level configuration object that enables and configures the supported
// transport stacks (currently HTTP and gRPC). Each nested transport config is optional:
// when nil or disabled, the corresponding transport constructors typically return nil
// and no server/client is created.
//
// The struct tags are compatible with the repository's config decoder (YAML/JSON/TOML).
type Config struct {
	// GRPC configures the gRPC transport stack (client and server).
	//
	// When nil or when the nested config is disabled, gRPC transport wiring is effectively
	// turned off and constructors such as `transport/grpc.NewServer` typically return nil.
	GRPC *grpc.Config `yaml:"grpc,omitempty" json:"grpc,omitempty" toml:"grpc,omitempty"`

	// HTTP configures the HTTP transport stack (client and server).
	//
	// When nil or when the nested config is disabled, HTTP transport wiring is effectively
	// turned off and constructors such as `transport/http.NewServer` typically return nil.
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
